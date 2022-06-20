package pkg

import (
	"sync"
	"testing"
)

func TestCreate(t *testing.T) {
	seq := NewSequence()
	if seq.Get() != 0 {
		t.Fatal("seq not initilized !")
	}
}

func TestSetandGet(t *testing.T) {
	seq := NewSequence()
	seq.Set(10)
	if seq.Get() != 10 {
		t.Fatal("set/get doesn't work !")
	}
}

func TestConcurrentGetAndSet(t *testing.T) {
	seq := NewSequence()
	wg := sync.WaitGroup{}
	add1 := func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			for {
				now := seq.Get()
				if seq.CompareAndSet(now, now+1) {
					break
				}
			}
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go add1()
	}

	wg.Wait()

	if seq.Get() != 100 {
		t.Fatal("concurrent get/set fail !")
	}
}

func TestCompareAndSet(t *testing.T) {
	seq := NewSequence8()
	seq.Set(10)
	fail := seq.CompareAndSet(20, 30)
	if fail == true {
		t.Fatal("compare and set failure case fail !")
	}
	succ := seq.CompareAndSet(10, 20)
	if succ == false {
		t.Fatal("compare and set success case fail !")
	}
	result := seq.Get()
	if result != 20 {
		t.Fatal("compare and set set case fail !")
	}
}

func TestConcurrentGetAndSet8(t *testing.T) {
	seq := NewSequence8()
	wg := sync.WaitGroup{}
	add1 := func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			for {
				now := seq.Get()
				if seq.CompareAndSet(now, now+1) {
					break
				}
			}
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go add1()
	}

	wg.Wait()

	if seq.Get() != 100 {
		t.Fatal("concurrent get/set fail !")
	}
}

func TestGetSeq(t *testing.T) {
	seq := int64(0)
	if GetSeq(&seq) != 0 {
		t.Fatal("seq not initilized !")
	}
}

func TestSetSeq(t *testing.T) {
	seq := int64(0)
	SetSeq(&seq, 10)
	if GetSeq(&seq) != 10 {
		t.Fatal("unable to set seq directly !")
	}
}
