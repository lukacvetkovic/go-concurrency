package main

import (
	"fmt"
	"sync"
)

func fanOut(messagesChannel <-chan int, processesNumber int) []chan int {
	messagesChannelArr := make([]chan int, processesNumber)
	for i := range messagesChannelArr {
		messagesChannelArr[i] = make(chan int, 0)
	}

	go func() {
		defer func() {
			for _, channel := range messagesChannelArr {
				close(channel)
			}
		}()
		for message := range messagesChannel {
			messagesChannelArr[message%processesNumber] <- message
		}
	}()

	return messagesChannelArr
}

func fanIn(messagesChannelArr ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(messagesChannelArr))

	out := make(chan int)

	for _, messageChannel := range messagesChannelArr {
		go func(messageChannel <-chan int) {
			defer wg.Done()

			for message := range messageChannel {
				out <- message
			}
		}(messageChannel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

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
	channelsForProcessing := fanOut(consumedMessages, processNumber)

	processedMessagesChannels := make([]<-chan int, 0, processNumber)

	for i := 0; i < processNumber; i++ {
		processed := process(channelsForProcessing[i])
		processedMessagesChannels = append(processedMessagesChannels, processed)
	}
	/////////////////////////////////////////////////////////

	aggregate := fanIn(processedMessagesChannels...)

	for message := range aggregate {
		fmt.Println(message)
	}

}
