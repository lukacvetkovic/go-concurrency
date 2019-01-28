package main

import (
	"sync"
)

func FanOut(messagesChannel <-chan int, processesNumber int) []chan int {
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

func FanIn(messagesChannelArr ...<-chan int) <-chan int {
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
