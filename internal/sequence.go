package pkg

import (
	"sync/atomic"
)

type Sequence struct {
	value [8]int64
}

func NewSequence() Sequence {
	return Sequence{value: [8]int64{0, 0, 0, 0, 0, 0, 0, 0}}
}

func (seq Sequence) Get() int64 {
	return atomic.LoadInt64(&seq.value[0])
}

func (seq *Sequence) Set(value int64) {
	atomic.StoreInt64(&seq.value[0], value)
}

func (seq *Sequence) CompareAndSet(original int64, value int64) bool {
	return atomic.CompareAndSwapInt64(&seq.value[0], original, value)
}
