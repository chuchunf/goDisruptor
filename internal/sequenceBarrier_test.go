package pkg

import (
	"testing"
)

func TestCreateNewBarrier(t *testing.T) {
	seq := NewSequence()
	wait := BusySpinWaitStrategy{}
	barrier := NewSequenceBarrier(wait, &seq)
	if barrier.cursor == nil {
		t.Fatal("unable to create new barrier !")
	}
}

func TestWaitFor(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	wait := BusySpinWaitStrategy{}
	barrier := NewSequenceBarrier(wait, &seq)
	result := barrier.WaitFor(99)
	if result != 99 {
		t.Fatal("barrier not able to wait for sequence !")
	}
}
