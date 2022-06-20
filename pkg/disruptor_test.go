package pkg

import (
	"sync"
	"testing"
	"time"
)

func TestCreateDisrutpor(t *testing.T) {
	disruptor := NewDisruptor[int64](1024)
	if disruptor == nil {
		t.Fatal("unable to create new disruptor")
	}
}

func TestCreateDisruptorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("no panic due to negative size")
		}
	}()
	NewDisruptor[int64](-1)
}

func TestAddnewProducer(t *testing.T) {
	disruptor := NewDisruptor[int64](1024)
	producer := disruptor.AddProducer(func(pooled *int64, updated int64) {})
	if producer == nil {
		t.Fatal("unable to add producer")
	}
}

func TestAddnewConsumer(t *testing.T) {
	disruptor := NewDisruptor[int64](1024)
	consumer := disruptor.AddConsumer(func(event *int64) {})
	if consumer == nil {
		t.Fatal("unable to add consumer")
	}
}

func TestDisruptor(t *testing.T) {
	disruptor := NewDisruptor[int64](1024)
	producer := disruptor.AddProducer(func(pooled *int64, updated int64) {
		*pooled = updated
	})
	consumer := disruptor.AddConsumer(func(event *int64) {
		if *event != 100 {
			t.Fatal("not able to process next event")
		}
	})

	producer(100)
	consumer()
}

func TestDisruptorConcurrently(t *testing.T) {
	disruptor := NewDisruptor[int64](1024)
	producer := disruptor.AddProducer(func(pooled *int64, updated int64) {
		*pooled = updated
	})
	consumer := disruptor.AddConsumer(func(event *int64) {
		if *event != 100 {
			t.Fatal("not able to process next event")
		}
	})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 2000; i++ {
			time.Sleep(1 * time.Millisecond)
			producer(100)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 2000; i++ {
			consumer()
		}
	}()

	wg.Wait()
}
