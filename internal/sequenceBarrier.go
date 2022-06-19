package pkg

type SequenceBarrier struct {
	waitStrategy WaitStrategy
	cursor       *Sequence
}

func NewSequenceBarrier(waitStrategy WaitStrategy, cursor *Sequence) SequenceBarrier {
	return SequenceBarrier{
		waitStrategy: waitStrategy,
		cursor:       cursor,
	}
}

func (barrier *SequenceBarrier) WaitFor(seq int64) int64 {
	return barrier.waitStrategy.waitFor(seq, barrier.cursor)
}
