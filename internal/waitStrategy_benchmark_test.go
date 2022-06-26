package pkg

import (
	"runtime/debug"
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

// from cpuprofile, the none GC version has simpler call graph and in general fast
// no GC version should be used for performance tuning
func BenchmarkBusySpinWithoutGC(b *testing.B) {
	debug.SetGCPercent(-1)

	seq := NewSequence()
	seq.Set(100)
	strategy := BusySpinWaitStrategy{}
	for i := 0; i < b.N; i++ {
		strategy.waitFor(99, &seq)
	}
}
