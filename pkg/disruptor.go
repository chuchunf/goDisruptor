package pkg

import (
	internal "goDisruptor/internal"
)

/*
** use reference so no mallocgc will be called, consumer should not modify the data
 */
type Consumer[E any] func(event *E)

type Producer[E any] func(pooled *E, updated E)

type Disrutpor[E any] struct {
	ringbuffer *internal.RingBuffer[E]
	consumers  []Consumer[E]
	barrier    internal.SequenceBarrier
}

func NewDisruptor[E any](size int64) *Disrutpor[E] {
	ring, err := internal.NewRingBuffer[E](size, internal.NewSequencer(size))
	if err != nil {
		panic(err)
	}

	return &Disrutpor[E]{
		ringbuffer: ring,
		barrier:    ring.CreateBarrier(),
	}
}

func (disruptor Disrutpor[E]) AddProducer(producer Producer[E]) func(event E) {
	return func(event E) {
		seq := disruptor.ringbuffer.Next()
		pooled := disruptor.ringbuffer.Get(seq)
		producer(pooled, event)
		disruptor.ringbuffer.Publish(seq)
	}
}

func (disruptor Disrutpor[E]) AddConsumer(consumer Consumer[E]) func() {
	seq := internal.NewSequence()
	disruptor.ringbuffer.AddGatingSequence(&seq)
	disruptor.consumers = append(disruptor.consumers, consumer)
	return func() {
		next := seq.Get() + 1
		next = disruptor.barrier.WaitFor(next)
		pooled := disruptor.ringbuffer.Get(next)
		consumer(pooled)
		seq.Set(next)
	}
}
