package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	go func() {
		for {
			t := rand.Intn(30)

			if t <= 30 && t >= 20 {
				ch <- t
			} else if t < 20 && t >= 10 {
				close(ch)
				return
			} else {
				continue
			}
		}
	}()

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-ch:
					if !ok {
						fmt.Println("Channel is closed")
						return
					}

					fmt.Println(data)
				default:
					//fmt.Println("DEFAULT")
				}
			}
		}()
	}

	wg.Wait()

	fmt.Println("OK!")
}
