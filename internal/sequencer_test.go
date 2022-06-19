package pkg

import (
	"testing"
	"time"
)

func TestCreateSequencer(t *testing.T) {
	seqcer := NewSequencer(1024)
	if seqcer.publishN(1, 10) != nil {
		t.Fatal("unable to create sequence and publish")
	}
}

func TestNextN(t *testing.T) {
	seqcer := NewSequencer(1024)
	seq1 := NewSequence()
	seq2 := NewSequence()
	seqcer.addGatingSequences(&seq1)
	seqcer.addGatingSequences(&seq2)
	seqcer.publish(10)
	next, _ := seqcer.next()
	if next != 11 {
		t.Fatal("sequence fail to claim next 1 !")
	}
}

func TestNextNError(t *testing.T) {
	seqcer := NewSequencer(1024)
	seq1 := NewSequence()
	seqcer.addGatingSequences(&seq1)
	seqcer.publish(10)
	_, err := seqcer.nextN(-1)
	if err != errorIllegalSizeRequired {
		t.Fatal("unable to handle negative nextN")
	}
	_, err = seqcer.nextN(1025)
	if err != errorIllegalSizeRequired {
		t.Fatal("unable to handle too big nextN")
	}
}

func TestNextNWrapped(t *testing.T) {
	seqcer := NewSequencer(1024)
	seqcer.publish(1100)

	seq1 := NewSequence()
	seq2 := NewSequence()
	seq1.Set(1000)
	seq2.Set(1000)
	seqcer.addGatingSequences(&seq1)
	seqcer.addGatingSequences(&seq2)
	next, _ := seqcer.nextN(100)
	if next != 1200 {
		t.Fatal("sequence fail to claim next 1 !")
	}
}

func TestNextConcurrentAccess(t *testing.T) {
	seqcer := NewSequencer(1024)
	seqcer.publish(1124)

	seq1 := NewSequence()
	seq1.Set(101)
	seqcer.addGatingSequences(&seq1)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1000 * time.Millisecond)
			now := seq1.Get()
			seq1.Set(now + 1)
		}
	}()

	for i := 0; i < 10; i++ {
		result, _ := seqcer.next()
		if int64(1124+i+1) != result {
			t.Fatal("failed to claim next 1 !")
		}
		seqcer.publish(result)
	}
}
