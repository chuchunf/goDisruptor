package example

import (
	"runtime"
	"sync"
	"testing"

	unix "golang.org/x/sys/unix"
)

func BenchmarkTickerPlant(b *testing.B) {
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

/*
** performance is worse with pinned CPU .. :( ..
** from call graphy, this pincpu version should be faster ...
 */
func BenchmarkTickerPlantPinCPU(b *testing.B) {
	if runtime.NumCPU() < 4 {
		panic("need 4 CPU cores for benchmark !")
	}
	var cpuset1, cpuset2, cpuset3 = unix.CPUSet{}, unix.CPUSet{}, unix.CPUSet{}
	cpuset1.Set(1)
	cpuset2.Set(2)
	cpuset3.Set(3)

	producer, logger, processer := TickerPlant()
	iterations := int(b.N)

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
