package pkg

import (
	"runtime"
	"time"
)

/*
** How consumers should wait for the next sequence to be available by producer
 */
type WaitStrategy interface {
	waitFor(next int64, seq *Sequence) int64
}

/*
** busy spin for the next available sequence
 */
type BusySpinWaitStrategy struct {
}

func (BusySpinWaitStrategy) waitFor(next int64, seq *Sequence) int64 {
	for next > seq.Get() {
	}
	return next
}

/*
** yield current execution
 */
type YieldWaitStrategy struct {
}

func (YieldWaitStrategy) waitFor(next int64, seq *Sequence) int64 {
	for next > seq.Get() {
		runtime.Gosched()
	}
	return next
}

/*
** sleep for 1 nano second
 */
type SleepWaitStrategy struct {
}

func (SleepWaitStrategy) waitFor(next int64, seq *Sequence) int64 {
	for next > seq.Get() {
		time.Sleep(1 * time.Nanosecond)
	}
	return next
}
