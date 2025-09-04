package pkg

import (
	"runtime"
	"sync"
	"testing"

	unix "golang.org/x/sys/unix"
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

func BenchmarkRingBufferPinCPU(b *testing.B) {
	if runtime.NumCPU() < 4 {
		panic("need 4 CPU cores for benchmark !")
	}
	var cpuset1, cpuset2, cpuset3 = unix.CPUSet{}, unix.CPUSet{}, unix.CPUSet{}
	cpuset1.Set(1)
	cpuset2.Set(2)
	cpuset3.Set(3)

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
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset1)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ring.Get(int64(i))
			seq1.Set(int64(i))
		}
	}
	reader2 := func() {
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset2)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ring.Get(int64(i))
			seq2.Set(int64(i))
		}
	}
	writer := func() {
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset3)
		defer runtime.UnlockOSThread()
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

func switchCPU(changeCPU bool, cpu1 unix.CPUSet, cpu2 unix.CPUSet) bool {
	runtime.UnlockOSThread()
	runtime.LockOSThread()
	if changeCPU {
		unix.SchedSetaffinity(0, &cpu1)
	} else {
		unix.SchedSetaffinity(0, &cpu2)
	}
	return !changeCPU
}

func BenchmarkRingBufferForceSwitchCPU(b *testing.B) {
	if runtime.NumCPU() < 8 {
		panic("need 8 CPU cores for benchmark !")
	}
	var cpuset1, cpuset2, cpuset3 = unix.CPUSet{}, unix.CPUSet{}, unix.CPUSet{}
	cpuset1.Set(1)
	cpuset2.Set(2)
	cpuset3.Set(3)
	var cpuset4, cpuset5, cpuset6 = unix.CPUSet{}, unix.CPUSet{}, unix.CPUSet{}
	cpuset1.Set(4)
	cpuset2.Set(5)
	cpuset3.Set(6)

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
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset1)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset1, cpuset4)
			}
			ring.Get(int64(i))
			seq1.Set(int64(i))
		}
	}
	reader2 := func() {
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset2)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset2, cpuset5)
			}
			ring.Get(int64(i))
			seq2.Set(int64(i))
		}
	}
	writer := func() {
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset3)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset3, cpuset6)
			}
			ring.Get(int64(i))
			ring.Publish(int64(i))
		}
	}

	go reader1()
	go reader2()
	go writer()

	wg.Wait()
}
