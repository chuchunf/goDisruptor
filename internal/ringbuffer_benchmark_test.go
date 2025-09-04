package pkg

import (
	"sync"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	for i := 0; i < b.N; i++ {
		ring.Get(0)
	}
}

func BenchmarkPublish(b *testing.B) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	for i := 0; i < b.N; i++ {
		ring.Publish(0)
	}
}

func BenchmarkNext(b *testing.B) {
	seqcer := NewSequencer(1024)
	seqcer.publish(2)
	seq1 := NewSequence()
	seq1.Set(0)
	seqcer.addGatingSequences(&seq1)

	ring, _ := NewRingBuffer[int64](1024, seqcer)
	for i := 0; i < b.N; i++ {
		ring.Next()
		seq1.Set(int64(i))
	}
}

func BenchmarkRingBuffer(b *testing.B) {
	iterations := b.N
	seqcer := NewSequencer(1024)

	seq1 := NewSequence()
	seq2 := NewSequence()

	seqcer.addGatingSequences(&seq1)
	seqcer.addGatingSequences(&seq2)

	ring, _ := NewRingBuffer[int64](1024, seqcer)

	wg := sync.WaitGroup{}
	wg.Add(3)

	reader1 := func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ring.Get(int64(i))
			seq1.Set(int64(i))
		}
	}
	reader2 := func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ring.Get(int64(i))
			seq2.Set(int64(i))
		}
	}
	writer := func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ring.Get(int64(i))
			ring.Publish(int64(i))
		}
	}

	go reader1()
	go reader2()
	go writer()

	wg.Wait()
}
