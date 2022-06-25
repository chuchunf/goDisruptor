package pkg

import (
	"runtime/debug"
	"testing"
)

func BenchmarkSequenceGet(b *testing.B) {
	seq := NewSequence()
	for i := 0; i < b.N; i++ {
		seq.Get()
	}
}

// slower without GC
func BenchmarkSequenceGetWithoutGC(b *testing.B) {
	debug.SetGCPercent(-1)

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

func BenchmarkCompareAndSet(b *testing.B) {
	seq := NewSequence8()
	now := seq.Get()
	next := now + 1
	for i := 0; i < b.N; i++ {
		seq.CompareAndSet(now, next)
		now = next
		next = now + 1
	}
}

func BenchmarkSequence8Get(b *testing.B) {
	seq := NewSequence8()
	for i := 0; i < b.N; i++ {
		seq.Get()
	}
}

func BenchmarkSequence8Set(b *testing.B) {
	seq := NewSequence8()
	for i := 0; i < b.N; i++ {
		seq.Set(int64(i))
	}
}

func BenchmarkGetSeq(b *testing.B) {
	seq := int64(0)
	for i := 0; i < b.N; i++ {
		GetSeq(&seq)
	}
}

func BenchmarkSetSeq(b *testing.B) {
	seq := int64(0)
	for i := 0; i < b.N; i++ {
		SetSeq(&seq, int64(i))
	}
}
