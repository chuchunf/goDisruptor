package pkg

import (
	"testing"
	"time"
)

func TestBusySpinWaitStrategy(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := BusySpinWaitStrategy{}
	result := strategy.waitFor(99, &seq)
	if result != 99 {
		t.Fatal("busy spin wait not working !")
	}
}

func TestBusySpinWaitStrategyWait(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := BusySpinWaitStrategy{}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		seq.Set(102)
	}()

	result := strategy.waitFor(101, &seq)
	if result != 101 {
		t.Fatal("busy spin wait not working !")
	}
}

func TestYieldWaitStrategy(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := YieldWaitStrategy{}
	result := strategy.waitFor(99, &seq)
	if result != 99 {
		t.Fatal("yield wait not working !")
	}
}

func TestYieldWaitStrategyWait(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := YieldWaitStrategy{}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		seq.Set(102)
	}()

	result := strategy.waitFor(101, &seq)
	if result != 101 {
		t.Fatal("yield wait not working !")
	}
}

func TestSleepWaitStrategy(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := SleepWaitStrategy{}
	result := strategy.waitFor(99, &seq)
	if result != 99 {
		t.Fatal("yield wait not working !")
	}
}

func TestSleepWaitStrategyWait(t *testing.T) {
	seq := NewSequence()
	seq.Set(100)
	strategy := SleepWaitStrategy{}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		seq.Set(102)
	}()

	result := strategy.waitFor(101, &seq)
	if result != 101 {
		t.Fatal("sleep wait not working !")
	}
}
