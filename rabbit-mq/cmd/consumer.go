package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	quueName = "q1"
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
		quueName, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
			time.Sleep(5 * time.Second)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
