package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	fmt.Println("CPUs\t\t", runtime.NumCPU())
	fmt.Println("Gorutines\t", runtime.NumGoroutine())

	counter := 0
	var wg sync.WaitGroup

	const gs = 100
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			v := counter
			runtime.Gosched() // Gosched yields the processor, allowing other goroutines to run.
			v++
			counter = v
			wg.Done()
		}()
		fmt.Println("Gorutines\t", runtime.NumGoroutine())
	}

	wg.Wait()

	fmt.Println("Gorutines\t", runtime.NumGoroutine())
	fmt.Println("Counter: ", counter)

}

//go run -race
