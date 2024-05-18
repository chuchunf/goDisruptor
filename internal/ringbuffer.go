package pkg

import (
	"errors"
)

type RingBuffer[E any] struct {
	size      int64
	indexMask int64
	entries   []E
	sequencer Sequencer
}

func NewRingBuffer[E any](size int64, sequencer Sequencer) (*RingBuffer[E], error) {
	if size < 1 {
		return nil, errorBufferSizeLessthan1
	}
	if size&(size-1) != 0 {
		return nil, errorBufferSizePowerof2
	}

	return &RingBuffer[E]{
		size:      size,
		indexMask: size - 1,
		entries:   make([]E, size),
		sequencer: sequencer,
	}, nil
}

func (ring *RingBuffer[E]) CreateBarrier() SequenceBarrier {
	return ring.sequencer.createBarrier()
}

func (ring *RingBuffer[E]) AddGatingSequence(seq *Sequence) {
	ring.sequencer.addGatingSequences(seq)
}

func (ring *RingBuffer[E]) Next() int64 {
	next, _ := ring.sequencer.next()
	return next
}

// TODO: performance tesitng for remainder operation with mask
func (ring *RingBuffer[E]) Get(index int64) *E {
	return &ring.entries[index&ring.indexMask]
}

func (ring *RingBuffer[E]) Publish(index int64) {
	ring.sequencer.publish(index)
}

var (
	errorBufferSizeLessthan1 = errors.New("buffer size must be at least 1")
	errorBufferSizePowerof2  = errors.New("buffer size must be power of 2")
)
