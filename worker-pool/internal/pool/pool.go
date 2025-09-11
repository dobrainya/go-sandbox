package pool

import (
	"context"
	"fmt"
	"math/rand"
	"pool-example/internal/pool/wp"
	"time"
)

var maxMessages int = 100
var messagesCounter = 0
var workersCount int = 10

type IPool interface {
	Create()
	Handle(m Message)
	Wait()
	Stats()
}

type Message struct {
	Value string
	Key   string
}

func getMessages() []Message {
	messagesCount := rand.Intn(maxMessages)
	messages := make([]Message, 0, messagesCount)

	for range messagesCount {
		messagesCounter++
		messages = append(messages, Message{"hello" + string(messagesCounter)})
	}

	return messages
}

func Process(dataFn func() []Message, handler func(workerId int, m Message)) {
	wp.InitWorkers(workersCount)
	var pool IPool = wp.New(handler)
	var totalMsgCnt int

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
l:
	for {
		select {
		case <-ctx.Done():
			break l
		default:
		}
		messages := dataFn()

		totalMsgCnt += len(messages)

		pool.Create()

		for _, message := range messages {
			pool.Handle(message)
		}

		pool.Wait()
	}

	pool.Stats()

	fmt.Println("Total messages:", totalMsgCnt)
}
