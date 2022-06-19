package example

import (
	"sync"
	"testing"
	"time"
)

const cycles = 10_000

func TestTickerPlant(t *testing.T) {
	producer, logger, processer := TickerPlant()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < cycles; i++ {
			time.Sleep(1 * time.Millisecond)
			producer(TickData{seq: int32(i)})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < cycles; i++ {
			logger()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < cycles; i++ {
			processer()
		}
	}()

	wg.Wait()
}
