package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"flag"

	"github.com/joho/godotenv"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exchangeName = "qex" // ""  - is default exange in rabbit
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	var msg string

	flag.StringVar(&msg, "m", "", "Message to send to message broker")

	flag.Parse()

	if len(msg) == 0 {
		fmt.Println("message was not provided. Use -h for more info")
		return
	}

	connectionString := os.Getenv("RABBITMQ_CONNECTION_URI")
	rmq_tls_certfile := os.Getenv("CLIENT_TLS_CERT_FILE")
	rmq_tls_keyfile := os.Getenv("CLIENT_TLS_KEY_FILE")
	rmq_tls_cafile := os.Getenv("CLIENT_TLS_CA_FILE")

	fmt.Println(connectionString)

	cfg := tls.Config{}
	cfg.RootCAs = x509.NewCertPool()

	caCert, err := os.ReadFile(rmq_tls_cafile)
	failOnError(err, "Unable to read CA bundle")
	cfg.RootCAs.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(rmq_tls_certfile, rmq_tls_keyfile)
	failOnError(err, "Unable to read certificate or key")
	cfg.Certificates = append(cfg.Certificates, cert)
	cfg.MinVersion = tls.VersionTLS12

	conn, err := amqp.DialTLS_ExternalAuth(connectionString, &cfg)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"q1",  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"qex",  // exchange name
		false,  // no-wait
		nil,    // arguments
	)

	failOnError(err, "Cannot bind queue to exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)
}
