package pkg

import (
	"runtime"
	"testing"
)

func BenchmarkNextN(b *testing.B) {
	runtime.SetCPUProfileRate(10000)

	seqcer := NewSequencer(1024)
	seqcer.publish(10)

	seq1 := NewSequence()
	seq1.Set(9)
	seqcer.addGatingSequences(&seq1)

	for i := 0; i < b.N; i++ {
		next, _ := seqcer.next()
		seqcer.publish(next)
		seq1.Set(next - 2)
	}
}

func BenchmarkMinSeq(b *testing.B) {
	runtime.SetCPUProfileRate(10000)

	seq1 := NewSequence()
	seq2 := NewSequence()
	seq1.Set(10)
	seq2.Set(20)
	seqs := []*Sequence{&seq1, &seq2}

	for i := 0; i < b.N; i++ {
		getMinSeq(seqs)
	}
}
