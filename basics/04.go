package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		ch := make(chan int)
		go myFunction("Hello world", ch, i)
		fmt.Println(<-ch)
	}

	ch1 := make(chan int)
	go fnc1(ch1)

	for a := range ch1 {
		fmt.Println(a)
	}
}

func myFunction(value string, ch chan int, i int) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(value, i)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	ch <- rand.Intn(100)
}

func fnc1(ch chan int) {
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 10; i++ {
		ch <- rand.Intn(100)
	}

	close(ch)
}
