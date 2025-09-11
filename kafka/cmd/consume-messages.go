package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	timeout = 1800
	groupId = "consumers-go"
	topics  = [...]string{"notifications", "sms", "push"}
)

func main() {
	consumers := make([]*kafka.Consumer, 0, len(topics))
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for range len(topics) {
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			// User-specific properties that you must set
			"bootstrap.servers":                  "localhost:9091,localhost:9092,localhost:9093",
			"socket.connection.setup.timeout.ms": 5000,

			// Fixed properties
			//"security.protocol": "PLAINTEXT",
			//"sasl.mechanisms":   "PLAIN",
			"group.id":          groupId,
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			fmt.Printf("Failed to create consumer. %s", err)
			continue
		}

		topic := topics[rand.Intn(len(topics))]
		err = c.SubscribeTopics([]string{topic}, nil)
		if err != nil {
			c.Close()
			fmt.Printf("Consumer cannot subscribes to topic. %s", err)
			continue
		}

		consumers = append(consumers, c)
	}

	for consumerId, c := range consumers {
		wg.Add(1)
		go func() {
			// Set up a channel for handling Ctrl-C, etc
			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
			defer wg.Done()
			defer c.Close()

			// Process messages
			run := true
			for run {
				select {
				case <-ctx.Done():
					fmt.Printf("Consumer [%d] caught signal %v. Terminating\n", consumerId)
					run = false
				case sig := <-sigchan:
					fmt.Printf("Consumer [%d] caught timeout. Terminating\n", consumerId, sig)
					run = false
				default:
					ev, err := c.ReadMessage(100 * time.Millisecond)

					if err == nil {
						fmt.Printf("Consumer [%d] consumed event from topic %s: key = %-10s value = %s\n", consumerId,
							ev.TopicPartition, string(ev.Key), string(ev.Value))
					} else if !err.(kafka.Error).IsTimeout() {
						// The client will automatically try to recover from all errors.
						// Timeout is not considered an error because it is raised by
						// ReadMessage in absence of messages.
						fmt.Printf("Consumer [%d] caught error: %v (%v)\n", consumerId, err, ev)
						run = false
					}
				}
			}
		}()
	}

	wg.Wait()

	fmt.Println("_____________All consumers terminated_____________")
}
