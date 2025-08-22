package pkg

import (
	"sync/atomic"
)

/*
** Implementation of Sequence, which is the fundamental data structure for the entire package.
**
** Sequence is just a thread-safe counter, which for the consumer and producer to coordinate by getting the next number
** then process to the next slot from the ring buffer.
 */

// Sequence as a structure with a single int64 as counter
type Sequence struct {
	value int64
}

func NewSequence() Sequence {
	return Sequence{value: 0}
}

func (seq *Sequence) Get() int64 {
	return atomic.LoadInt64(&seq.value)
}

func (seq *Sequence) Set(value int64) {
	atomic.StoreInt64(&seq.value, value)
}

func (seq *Sequence) CompareAndSet(original int64, value int64) bool {
	return atomic.CompareAndSwapInt64(&seq.value, original, value)
}

// Sequence8 Sequence as a structure with a padded int64 array to avoid false sharing
type Sequence8 struct {
	value [8]int64
}

func NewSequence8() Sequence8 {
	return Sequence8{value: [8]int64{0, 0, 0, 0, 0, 0, 0, 0}}
}

func (seq *Sequence8) Get() int64 {
	return atomic.LoadInt64(&seq.value[0])
}

func (seq *Sequence8) Set(value int64) {
	atomic.StoreInt64(&seq.value[0], value)
}

func (seq *Sequence8) CompareAndSet(original int64, value int64) bool {
	return atomic.CompareAndSwapInt64(&seq.value[0], original, value)
}

// GetSeq uses an int64 directly
func GetSeq(seq *int64) int64 {
	return atomic.LoadInt64(seq)
}

func SetSeq(seq *int64, value int64) {
	atomic.StoreInt64(seq, value)
}
