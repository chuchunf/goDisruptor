package pkg

import (
	internal "goDisruptor/internal"
	"sync"
	"testing"
)

func BenchmarkDisruptor(b *testing.B) {
	size := int64(1024)
	ringbuffer, _ := internal.NewRingBuffer[int64](size, internal.NewSequencer(size))
	barrier := ringbuffer.CreateBarrier()

	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ReportAllocs()
	b.ResetTimer()

	iterations := int32(b.N)
	go func() {
		defer wg.Done()
		for i := int32(0); i < iterations; i++ {
			seq := ringbuffer.Next()
			pooled := ringbuffer.Get(seq)
			*pooled = int64(i)
			ringbuffer.Publish(seq)
		}
	}()

	go func() {
		defer wg.Done()
		seq := internal.NewSequence()
		ringbuffer.AddGatingSequence(&seq)
		for i := int32(0); i < iterations; i++ {
			next := seq.Get() + 1
			next = barrier.WaitFor(next)
			ringbuffer.Get(next)
			seq.Set(next)
		}
	}()

	wg.Wait()
}
