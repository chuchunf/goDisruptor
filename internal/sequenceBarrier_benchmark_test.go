package pkg

import (
	"testing"
)

func BenchmarkWaitFor(b *testing.B) {
	seq := NewSequence()
	seq.Set(100)
	wait := BusySpinWaitStrategy{}
	barrier := NewSequenceBarrier(wait, &seq)
	for i := 0; i < b.N; i++ {
		barrier.WaitFor(99)

	}
}
