package pkg

import (
	"testing"
)

func TestCreateRingBuffer(t *testing.T) {
	_, err := NewRingBuffer[int64](1024, NewSequencer(1024))
	if err != nil {
		t.Fatal("cannot create ring buffer !")
	}
}

func TestCreateRingBufferInvalidSize(t *testing.T) {
	_, err := NewRingBuffer[int64](-10, NewSequencer(1024))
	if err != errorBufferSizeLessthan1 {
		t.Fatal("cannot create size less than 1 !")
	}
}

func TestCreateRingBufferNotPowerOf2(t *testing.T) {
	_, err := NewRingBuffer[int64](99, NewSequencer(99))
	if err != errorBufferSizePowerof2 {
		t.Fatal("buffer size must be power of 2!")
	}
}

func TestRingBufferGet(t *testing.T) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	result := ring.Get(0)
	if *result != 0 {
		println("result is ", result)
		t.Fatal("not able to get element by index!")
	}
}

func TestRingBufferNext(t *testing.T) {
	seqcer := NewSequencer(1024)
	seqcer.publish(1100)
	seq1 := NewSequence()
	seq1.Set(1000)
	seqcer.addGatingSequences(&seq1)

	ring, _ := NewRingBuffer[int64](1024, seqcer)
	result := ring.Next()
	if result != 1101 {
		t.Fatal("cannot claim next sequence !")
	}
}

func TestRingBufferPublish(t *testing.T) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	ring.Publish(1)
}

func TestRingBufferCreateBarrier(t *testing.T) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	ring.CreateBarrier()
}

func TestRingBufferAddGatingSequence(t *testing.T) {
	ring, _ := NewRingBuffer[int64](1024, NewSequencer(1024))
	seq1 := NewSequence()
	ring.AddGatingSequence(&seq1)
}
