package pkg

import (
	"testing"
)

func BenchmarkBusySpin(b *testing.B) {
	seq := NewSequence()
	seq.Set(100)
	strategy := BusySpinWaitStrategy{}
	for i := 0; i < b.N; i++ {
		strategy.waitFor(99, &seq)
	}
}
