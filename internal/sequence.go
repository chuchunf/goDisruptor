package pkg

import (
	"sync/atomic"
)

/*
** Implementation of Sequence, which is the fundemtal data strcture for the entire package.
**
** Sequence is just a thread-safe counter, which for the consumer and producer to determine, get then process
** a slot from the ring buffer.
 */

// Sequence as a structure with a single int64 as counter
type Sequence struct {
	value int64
}

func NewSequence() Sequence {
	return Sequence{value: 0}
}

func (seq Sequence) Get() int64 {
	return atomic.LoadInt64(&seq.value)
}

func (seq *Sequence) Set(value int64) {
	atomic.StoreInt64(&seq.value, value)
}

func (seq *Sequence) CompareAndSet(original int64, value int64) bool {
	return atomic.CompareAndSwapInt64(&seq.value, original, value)
}

// Sequence as a structure with a paded int64 array as
type Sequence8 struct {
	value [8]int64
}

func NewSequence8() Sequence8 {
	return Sequence8{value: [8]int64{0, 0, 0, 0, 0, 0, 0, 0}}
}

func (seq Sequence8) Get() int64 {
	return atomic.LoadInt64(&seq.value[0])
}

func (seq *Sequence8) Set(value int64) {
	atomic.StoreInt64(&seq.value[0], value)
}

func (seq *Sequence8) CompareAndSet(original int64, value int64) bool {
	return atomic.CompareAndSwapInt64(&seq.value[0], original, value)
}

// Use an int64 directly
func GetSeq(seq *int64) int64 {
	return atomic.LoadInt64(seq)
}

func SetSeq(seq *int64, value int64) {
	atomic.StoreInt64(seq, value)
}
