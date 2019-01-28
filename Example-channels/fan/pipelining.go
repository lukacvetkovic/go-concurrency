package main

import (
	"fmt"
)

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

const processNumber int = 4

func main() {

	consumedMessages := consume()

	/////////////////////////////////////////////////////////
	channelsForProcessing := FanOut(consumedMessages, processNumber)

	processedMessagesChannels := make([]<-chan int, 0, processNumber)

	for i := 0; i < processNumber; i++ {
		processed := process(channelsForProcessing[i])
		processedMessagesChannels = append(processedMessagesChannels, processed)
	}
	/////////////////////////////////////////////////////////

	aggregate := FanIn(processedMessagesChannels...)

	for message := range aggregate {
		fmt.Println(message)
	}

}
