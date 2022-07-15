package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// A() info -> channel
// B() channel -> info

func main() {
	wg := &sync.WaitGroup{}
	IdsChannel := make(chan string)
	fakeIdsChan := make(chan string)

	closeChannels := make(chan int)

	wg.Add(3)

	go generateIds(wg, IdsChannel, closeChannels)
	go generateFalseIds(wg, fakeIdsChan, closeChannels)
	go logIds(IdsChannel, wg, fakeIdsChan, closeChannels)

	wg.Wait()
}

func generateFalseIds(waitGroup *sync.WaitGroup, fakeIdsChan chan<- string, closeChannels chan<- int) {
	for i := 0; i < 50; i++ {
		id := uuid.New()
		fakeIdsChan <- fmt.Sprintf("%d. %s", i+1, id.String())
	}

	close(fakeIdsChan)
	closeChannels <- 1

	waitGroup.Done()
}

func generateIds(waitGroup *sync.WaitGroup, idsChan chan<- string, closeChannels chan<- int) {
	for i := 0; i < 100; i++ {
		id := uuid.New()
		idsChan <- fmt.Sprintf("%d. %s", i+1, id.String())
	}

	close(idsChan)
	closeChannels <- 1

	waitGroup.Done()
}

func logIds(idsChan <-chan string, wg *sync.WaitGroup, fakeIdsChan <-chan string, closeChannels chan int) {
	closedCounter := 0
	for {
		select {
		case id, ok := <-idsChan:
			if ok {
				fmt.Println("Id: ", id)
			}
		case id, ok := <-fakeIdsChan:
			if ok {
				fmt.Println("Fake id: ", id)
			}
		case count, ok := <-closeChannels:
			if ok {
				closedCounter += count
			}
		}

		if closedCounter == 2 {
			close(closeChannels)
			break
		}
	}
	wg.Done()
}
