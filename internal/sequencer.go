package pkg

import (
	"errors"
)

type Sequencer interface {
	next() (int64, error)
	nextN(n int64) (int64, error)
	publish(seq int64) error
	publishN(low int64, high int64) error
	addGatingSequences(seq *Sequence)
	createBarrier() SequenceBarrier
}

type SingleProducerSequencer struct {
	size           int64
	cursor         Sequence    // writing seq
	gatingSequence []*Sequence // reading seqs
}

func NewSequencer(bufferSize int64) Sequencer {
	return &SingleProducerSequencer{
		size:   bufferSize,
		cursor: NewSequence(),
	}
}

/*
** add new gate sequence from readers
 */
func (seqcer *SingleProducerSequencer) addGatingSequences(seq *Sequence) {
	seqcer.gatingSequence = append(seqcer.gatingSequence, seq)
}

/*
** get barrier for consumers
 */
func (seqcer *SingleProducerSequencer) createBarrier() SequenceBarrier {
	return NewSequenceBarrier(BusySpinWaitStrategy{}, &seqcer.cursor)
}

/*
** claim 1/n sequence for writing
 */
func (seqcer *SingleProducerSequencer) next() (int64, error) {
	return seqcer.nextN(1)
}

func (seqcer *SingleProducerSequencer) nextN(n int64) (int64, error) {
	if n < 1 || n >= seqcer.size {
		return 0, errorIllegalSizeRequired
	}

	written := seqcer.cursor.Get()
	next := written + n
	wrap := next - seqcer.size

	for minValue := getMinSeq(seqcer.gatingSequence); wrap >= minValue; {
		minValue = getMinSeq(seqcer.gatingSequence)
	}

	return next, nil
}

func getMinSeq(seqs []*Sequence) int64 {
	var minValue int64
	for i, v := range seqs {
		if i == 0 || v.Get() < minValue {
			minValue = v.Get()
		}
	}
	return minValue
}

/*
** once the event(s) is/are updated in the ring buffer, update the written cursor and notify all readers
 */
func (seqcer *SingleProducerSequencer) publish(seq int64) error {
	seqcer.cursor.Set(seq)
	return nil
}

func (seqcer *SingleProducerSequencer) publishN(low int64, high int64) error {
	return seqcer.publish(high)
}

var (
	errorIllegalSizeRequired = errors.New("the size required exceeds the limit")
)
