package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	async1()
	async2()
	async3()
	async4()
}

func async1() {
	var money atomic.Int32
	wg := sync.WaitGroup{}

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			money.Add(1)
		}()
	}

	wg.Wait()
	fmt.Println(money.Load())
}

func async2() {
	money := 0
	donations := 0
	lock := sync.RWMutex{}
	wg := sync.WaitGroup{}

	go func() {
		for {

			lock.RLock()
			m := money
			d := donations
			lock.RUnlock()

			if d != m {
				fmt.Println(" Money is not equal to donations")
				return
			}
		}
	}()

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			lock.Lock()
			money++
			donations++
			lock.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(money)
}

func randomDurationWorkingFn(duration int) int {
	t := rand.Intn(duration)
	time.Sleep(time.Duration(t) * time.Second)

	return t
}

func async3() {
	start := time.Now()
	total := 0
	locker := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			t := randomDurationWorkingFn(5)
			locker.Lock()
			total += t
			locker.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("duration: ", time.Since(start))
	fmt.Println("total: ", total)
}

func async4() {
	start := time.Now()
	total := 0
	ch := make(chan int)

	for range 100 {
		go func() {
			t := randomDurationWorkingFn(5)
			ch <- t
		}()
	}

	for range 100 {
		total += <-ch
	}

	fmt.Println("duration: ", time.Since(start))
	fmt.Println("total: ", total)
}
