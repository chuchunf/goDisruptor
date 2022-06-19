package pkg

import (
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
