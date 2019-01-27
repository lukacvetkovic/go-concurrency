package main

import "fmt"

func consume() <-chan int {

	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i < 100; i++ {
			out <- i
		}
	}()

	return out
}

func process(input <-chan int) <-chan int {

	out := make(chan int)
	go func() {
		defer close(out)
		for value := range input {
			out <- value * value
		}
	}()

	return out
}

func main() {

	consumedMessages := consume()

	processedMessages := process(consumedMessages)

	for message := range processedMessages {
		fmt.Println(message)
	}

}
