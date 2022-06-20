package pkg

import (
	"sync/atomic"
)

/*
** mallocgc allocates 8 bytes cost around 8ns, compare to 20ns for 64 bytes
** atomic.StoreInt64 doesn't call mallocgc, it is consistent 4ns
 */
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

/*
** seems go lang allocate struct in heap, [8]int64 cost around 20ns
 */
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

// implementation using int64 directly, no allocation to heap, 0.2 ns only
func GetSeq(seq *int64) int64 {
	return atomic.LoadInt64(seq)
}

func SetSeq(seq *int64, value int64) {
	atomic.StoreInt64(seq, value)
}
