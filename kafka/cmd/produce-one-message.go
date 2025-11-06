package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	timeout       = 1800
	flushTimeInMs = 15000
	users         = [...]string{"Kate", "Max", "John", "Jake", "Joose", "Rozi"}
	messages      = [...]string{"message1", "message2", "message3", "message4", "message5"}
	topic         = "notifications"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	workerId := rand.Intn(100)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers":                  "localhost:9091,localhost:9092,localhost:9093",
		"socket.connection.setup.timeout.ms": 1000,
		"acks":                               "all",
	})
	defer p.Close()

	if err != nil {
		fmt.Printf("Failed to create producer. %s\n", err)
		return
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigchan)

	ch := make(chan kafka.Event)
	defer close(ch)

	select {
	case sig := <-sigchan:
		fmt.Printf("Producer [%d] caught signal %v. Terminating\n", workerId, sig)
		return
	case <-ctx.Done():
		fmt.Printf("Producer [%d] caught timeout. Terminating\n", workerId)
		return
	default:
	}

	key := users[0]
	message := messages[rand.Intn(len(messages))]

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(message),
	}, ch)

	if err != nil {
		fmt.Printf(
			"Producer [%d] cannot produces message \"%s\" to the topic \"%s\" for key \"%s\": %s",
			workerId,
			message,
			topic,
			key,
			err.Error(),
		)

		return
	}

	event, ok := <-ch

	if !ok {
		fmt.Printf("Producer's [%d] channel is closed", workerId)
		return
	}

	switch ev := event.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			fmt.Printf("Producer [%d] is failed to deliver message: %v\n", workerId, ev.TopicPartition)

			return
		}

		fmt.Printf(
			"Producer [%d] produced event to topic %s: key = %-10s value = %s\n",
			workerId,
			ev.TopicPartition,
			string(ev.Key),
			string(ev.Value),
		)

		p.Flush(flushTimeInMs)
	}
}
