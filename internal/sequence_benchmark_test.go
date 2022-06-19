package pkg

import (
	"sync/atomic"
	"testing"
)

func BenchmarkCompareAndSet(b *testing.B) {
	seq := NewSequence()
	now := seq.Get()
	next := now + 1
	for i := 0; i < b.N; i++ {
		seq.CompareAndSet(now, next)
		now = next
		next = now + 1
	}
}

func BenchmarkSequenceGet2(b *testing.B) {
	seq := int64(100)
	for i := 0; i < b.N; i++ {
		atomic.LoadInt64(&seq)
	}
}

func BenchmarkSequenceGet(b *testing.B) {
	seq := NewSequence()
	for i := 0; i < b.N; i++ {
		seq.Get()
	}
}

func BenchmarkSequenceSet(b *testing.B) {
	seq := NewSequence()
	for i := 0; i < b.N; i++ {
		seq.Set(int64(i))
	}
}
