package example

import (
	"runtime"
	"runtime/debug"
	"sync"
	"testing"

	unix "golang.org/x/sys/unix"
)

func BenchmarkTickerPlant(b *testing.B) {
	debug.SetGCPercent(-1)

	producer, logger, processer := TickerPlant()
	iterations := int32(b.N)

	wg := sync.WaitGroup{}
	wg.Add(3)

	b.ReportAllocs()
	b.ResetTimer()

	ticker := TickData{seq: 0}

	go func() {
		defer wg.Done()
		for i := int32(0); i < iterations; i++ {
			ticker.seq = i
			producer(ticker)
		}
	}()

	go func() {
		defer wg.Done()
		for i := int32(0); i < iterations; i++ {
			logger()
		}
	}()

	go func() {
		defer wg.Done()
		for i := int32(0); i < iterations; i++ {
			processer()
		}
	}()

	wg.Wait()
}

func BenchmarkTickerPlantPinCPU(b *testing.B) {
	debug.SetGCPercent(-1)

	if runtime.NumCPU() < 4 {
		panic("need 4 CPU cores for benchmark !")
	}
	var cpuset1, cpuset2, cpuset3 = unix.CPUSet{}, unix.CPUSet{}, unix.CPUSet{}
	cpuset1.Set(1)
	cpuset2.Set(2)
	cpuset3.Set(3)

	producer, logger, processer := TickerPlant()
	iterations := b.N

	wg := sync.WaitGroup{}
	wg.Add(3)

	b.ReportAllocs()
	b.ResetTimer()

	ticker := TickData{seq: 0}

	go func() {
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset1)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ticker.seq = int32(i)
			producer(ticker)
		}
	}()

	go func() {
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset2)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			logger()
		}
	}()

	go func() {
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset3)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			processer()
		}
	}()

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

func BenchmarkTickerPlantPinCPUForceSwitch(b *testing.B) {
	debug.SetGCPercent(-1)

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

	producer, logger, processer := TickerPlant()
	iterations := b.N

	wg := sync.WaitGroup{}
	wg.Add(3)

	b.ReportAllocs()
	b.ResetTimer()

	ticker := TickData{seq: 0}

	go func() {
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset1)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset1, cpuset4)
			}
			ticker.seq = int32(i)
			producer(ticker)
		}
	}()

	go func() {
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset2)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset2, cpuset5)
			}
			logger()
		}
	}()

	go func() {
		changeCpu := true
		runtime.LockOSThread()
		unix.SchedSetaffinity(0, &cpuset3)
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%100_000 == 0 {
				changeCpu = switchCPU(changeCpu, cpuset3, cpuset6)
			}
			processer()
		}
	}()

	wg.Wait()
}
