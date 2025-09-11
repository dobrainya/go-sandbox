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
	workerSleepMaxTime = 5
	timeout            = 1800
	flushTimeInMs      = 15000
	users              = [...]string{"Kate", "Max", "John", "Jake", "Joose", "Rozi"}
	messages           = [...]string{"message1", "message2", "message3", "message4", "message5"}
	topics             = [...]string{"notifications", "sms", "push"}
)

func main() {
	producers := make([]*kafka.Producer, 0, len(topics))
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for id := range len(topics) {
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			// User-specific properties that you must set
			"bootstrap.servers":                  "localhost:9091,localhost:9092,localhost:9093",
			"socket.connection.setup.timeout.ms": 1000,
			"acks":                               "all",
		})
		_ = id
		if err != nil {
			fmt.Printf("Failed to create producer. %s\n", err)
			continue
		}

		// Go-routine to handle message delivery reports and
		// possibly other event types (errors, stats, etc)
		// go func() {
		// 	for e := range p.Events() {
		// 		switch ev := e.(type) {
		// 		case *kafka.Message:
		// 			if ev.TopicPartition.Error != nil {
		// 				fmt.Printf("Producer %d: Failed to deliver message: %v\n", id, ev.TopicPartition)
		// 			} else {
		// 				fmt.Printf("Producer %d: Produced event to topic %s: key = %-10s value = %s\n", id,
		// 					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		// 			}
		// 		}
		// 	}
		// }()
		producers = append(producers, p)
	}

	for id, p := range producers {
		wg.Add(1)
		go func() {
			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
			defer wg.Done()
			defer p.Close()
			ch := make(chan kafka.Event)
			//defer close(ch)
			for {
				select {
				case sig := <-sigchan:
					fmt.Printf("Producer [%d] caught signal %v. Terminating\n", id, sig)
					return
				case <-ctx.Done():
					fmt.Printf("Producer [%d] caught timeout. Terminating\n", id)
					return
				default:
				}

				// produceFunctionAsync(ctx, id, p)

				produceFunctionSync(id, p, ch)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("________________All poducers are terminated________________\n")
}

// func produceFunctionAsync(ctx context.Context, workerId int, p *kafka.Producer) {
// 	key := users[rand.Intn(len(users))]
// 	message := messages[rand.Intn(len(messages))]
// 	topic := topics[rand.Intn(len(topics))]
// 	err := p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Key:            []byte(key),
// 		Value:          []byte(message),
// 	}, nil)

// 	if err != nil {
// 		fmt.Printf(
// 			"Producer [%d] cannot produce message \"%s\" to the topic \"%s\" for key \"%s\": %s",
// 			workerId,
// 			message,
// 			topic,
// 			key,
// 			err.Error(),
// 		)

// 		return
// 	}

// 	p.Flush(flushTime * 1000)

// 	sleepTime := rand.Intn(workerSleepMaxTime)

// 	select {
// 	case <-ctx.Done():
// 		return
// 	case <-time.After(time.Duration(sleepTime) * time.Second):
// 	}
// }

func produceFunctionSync(workerId int, p *kafka.Producer, ch chan kafka.Event) {
	key := users[rand.Intn(len(users))]
	message := messages[rand.Intn(len(messages))]
	topic := topics[workerId]
	err := p.Produce(&kafka.Message{
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

	sleepTime := rand.Intn(workerSleepMaxTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)

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
