package pkg

import (
	"runtime/debug"
	"sync"
	"testing"
)

func BenchmarkSequenceGet(b *testing.B) {
	seq := NewSequence()
	for i := 0; i < b.N; i++ {
		seq.Get()
	}
}

// slower without GC with vscoder, but faster when triggered directly
// no significant difference in call graph, likely due to the I/O, lock, Timer, scheduling etc.
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

func BenchmarkSequenceSetWithoutGC(b *testing.B) {
	debug.SetGCPercent(-1)

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

func BenchmarkSequence8GetWithoutGC(b *testing.B) {
	debug.SetGCPercent(-1)

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

func BenchmarkConcurrentGetSetRaw(b *testing.B) {
	debug.SetGCPercent(-1)
	seq := int64(0)
	iterations := int64(b.N)
	wg := sync.WaitGroup{}
	wg.Add(3)

	add := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			SetSeq(&seq, i)
		}
	}

	get := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			GetSeq(&seq)
		}
	}

	go add()
	go get()
	go get()
	wg.Wait()
}

func BenchmarkConcurrentGetSetWithGC(b *testing.B) {
	seq := NewSequence()
	iterations := int64(b.N)
	wg := sync.WaitGroup{}
	wg.Add(3)

	add := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Set(i)
		}
	}

	get := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Get()
		}
	}

	go add()
	go get()
	go get()
	wg.Wait()
}

func BenchmarkConcurrentGetSet(b *testing.B) {
	debug.SetGCPercent(-1)

	seq := NewSequence()
	iterations := int64(b.N)
	wg := sync.WaitGroup{}
	wg.Add(3)

	add := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Set(i)
		}
	}

	get := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Get()
		}
	}

	go add()
	go get()
	go get()
	wg.Wait()
}

func BenchmarkConcurrentGetAndSet8(b *testing.B) {
	debug.SetGCPercent(-1)

	seq := NewSequence8()
	wg := sync.WaitGroup{}
	wg.Add(3)
	iterations := int64(b.N)

	add := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Set(i)
		}
	}

	get := func() {
		defer wg.Done()
		for i := int64(0); i < iterations; i++ {
			seq.Get()
		}
	}

	go add()
	go get()
	go get()
	wg.Wait()
}
