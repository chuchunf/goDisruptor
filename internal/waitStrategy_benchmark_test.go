package pkg

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"
)

func BenchmarkBusySpin(b *testing.B) {
	seq := NewSequence()
	seq.Set(100)
	strategy := BusySpinWaitStrategy{}
	for i := 0; i < b.N; i++ {
		strategy.waitFor(99, &seq)
	}
}

func BenchmarkBusySpin2(b *testing.B) {
	for i := 0; i < b.N; i++ {
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

func BenchmarkYieldWait(b *testing.B) {
	seq := NewSequence()
	seq.Set(100)
	strategy := YieldWaitStrategy{}
	for i := 0; i < b.N; i++ {
		strategy.waitFor(99, &seq)
	}
}

func TestBusySpin(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)

	go func() {
		time.Sleep(1 * time.Second)
		seq.Set(102)
	}()

	count := 0
	for 101 > seq.Get() {
		count = count + 1
		// busy spin
	}
	t.Logf("processed %d", count)
}

func TestSleep(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)

	go func() {
		time.Sleep(1 * time.Second)
		seq.Set(102)
	}()

	count := 0
	for 101 > seq.Get() {
		count = count + 1
		time.Sleep(1 * time.Nanosecond)
	}
	t.Logf("processed %d", count)
}

func TestYield(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)

	go func() {
		time.Sleep(1 * time.Second)
		seq.Set(102)
	}()

	count := 0
	for 101 > seq.Get() {
		count = count + 1
		runtime.Gosched()
	}
	t.Logf("processed %d", count)
}
