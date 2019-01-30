package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {

	fmt.Println("CPUs\t\t", runtime.NumCPU())
	fmt.Println("Gorutines\t", runtime.NumGoroutine())

	var counter int32 = 0
	var wg sync.WaitGroup

	const gs = 100
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			runtime.Gosched() // Gosched yields the processor, allowing other goroutines to run.
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Gorutines\t", runtime.NumGoroutine())
	fmt.Println("Counter: ", counter)

}

//go run -race
