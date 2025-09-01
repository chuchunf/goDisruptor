
# Introduction
The LMAX Disruptor is a high-performance inter-thread messaging library developed in Java, renowned for its exceptionally low latency  
by using lock-free algorithms, careful memory layout and other consideration at hardware level to minimize contention and improve the throughput.

The objective of the port is to learn to implement, test and verify low latency library in general, as golang prefers goroutine.
As Go provides powerful built-in primitives like channels and goroutines.

This document focuses on the performance analysis of this Go port of the LMAX Disruptor. 

$-$

# Design of LMAX-Disruptor
Main components from bottom up: 
* Sequence: A thread-safe, atomic counter used to track the progress of a Producer or Consumer. It is the foundation of the lock free coordination mechanism.
* WaitStrategy: Strategies or mechanisms for the producer or consumer to wait for the sequence.
* SequenceBarrier: A wrapper that uses a given WaitStrategy to protect write access of a given Sequence (for Consumer).
* Sequencer: The main implementation, 
  * maintain a cursor as a sequence (for the ringbuffer), 
  * a list of gating sequences for consumer 
  * a list of Barriers for producers 
  * provide core logic of coordinating the access of read and write with gating sequences and barriers 
    * next: find next available slot by comparing the next value of the cursor (of the ring buffer) or all gating sequences of all consumers, the next value must be larger that any of them to ensure the slot is not still used by any of the Consumer
    * publish: it is used by publisher to increate the counter and notify all consumers the slot is now available 
* RingerBuffer: maintain the actual ring buffer data structure in a memory aligned way for fast access and a sequencer to coordinate the access to the ring
* Disrupter: maintain the RingerBuffer and list of consumers



$-$

# Benchmark and Performance Analysis

## Sequence
Sequence is the fundamental data structure for the entire package.

In short, sequence is just a thread-safe counter, points to a slot in the ring buffer.
The consumer and producer uses a wait strategy to wait for the number to be available in the sequence
Once get the number/slot, the consumer or producer can continue its process.

Following Sequences implementation are tested
* a struct with int64
* a struct with [8]int64 (padding for 64 bytes to avoid false sharing)
* directly with int64

All performance testing are conducted in Windows 11 with Ryzen 5 7600 @ 3.8GHz in console

```bash
# benchmark 
go test -benchmem -run=^$ -tags -race -bench ^BenchmarkSequenceGet$ goDisruptor/internal

# profile
go test -cpuprofile cpu.prof -memprofile mem.prof -benchmem -run=^$ -tags -race -bench ^BenchmarkSequenceGet$ goDisruptor/internal

# get top function calls
go tool pprof cpu.prof => top

# generate call graph
go tool pprof cpu.prof => png

# get top memory allocation
go tool pprof mem.prof => top
```

$~$

### Benchmark results
#### for Get and Set
|Implemenation   | Get without GC | Get with GC  | Set without GC | Set with GC   | 
|----------------|----------------|--------------|---------------|---------------|
|struct int64   | 0.2182 ns/op   | 0.2119 ns/op | 1.610 ns/op   | 1.688 ns/op   | 
|struct [8]int64 | 0.2142 ns/op   | 0.2194 ns/op | 1.586 ns/op   | 1.607 ns/op   |
|int64          | 0.2119 ns/op   | 0.2121 ns/op | 1.612 ns/op   | 1.611 ns/op   |

#### for concurrent Get and Set
|Implementation                   | Get and Set  | Test code                                                        |
|---------------------------------|--------------|------------------------------------------------------------------|
|struct int64                    | 1.944 ns/op  | [BenchmarkConcurrentGetSetRaw](sequence_benchmark_test.go#L137)  |
|struct [8]int64                 | 1.958 ns/op  | [BenchmarkConcurrentGetAndSet8](sequence_benchmark_test.go#L221) |
|int64                           | 1.948 ns/op  | [BenchmarkConcurrentGetSet](sequence_benchmark_test.go#L192)     |
|False sharing int64             | **14.91 ns/op**  | [BenchmarkFalseSharing](sequence_benchmark_test.go#L250)     |
|No False Sharing struct [8]int64 | 1.931 ns/op  | [BenchmarkNoFalseSharing](sequence_benchmark_test.go#L295)       |

$~$

### Analysis
#### Impact of sampling rate
Given we're testing for nano second changes and the function used to run less than 1 second,
the default sample rate of 100 is too small.
For all the benchmarking, we will use sample rate at 10000
```go
    runtime.SetCPUProfileRate(10000)
```

$~$

#### Impact of cache line / false sharing
From the table for concurrent get and set, we can clearly tell that the false sharing have a significant impact
on performance. Additional note, struct[8]int64 dose provide the memory alignment which ensure the data fall into
the same cache line and avoid false sharing.

$~$

#### Impact of GC
From the benchmark above, the impact of GC in all the cases are negligible, probably due to the fact that
the memory footprint and change are relatively small, there are not much GC activities anyway.

This can be confirmed by the cpu.prof that the call graph are almost identical.
Although it needs to be noted that from the mem.prof, the tests without GC allocate a larger memory initially.

$~$

#### Impact of malloc
Impact on malloc is also negligible, probably due to the memory allocation is quite small.

$~$

#### Impact of the struct
The impact on additional struct is negligible too, the call graphy and memory profiling are almost identical.
Golang compiler has optimization on this already.

$-$

## WaitStrategy
WaitStrategy defines the strategy or mechanism for the producer or consumer to wait for the sequence.

Following waiting strategies are implemented and tested
* BusySpinWaitStrategy: waiting in an infinite loop
* YieldWaitStrategy: yield the current thread's execution (via runtime.Gosched())
* SleepWaitStrategy: let current thread sleep for 1 Nano seconds (via time.Sleep())

$-$

### Benchmark results
Obviously, the BusySpin strategy is the best performance implementation which has the best through put at a cost of wasting CPU cycle.
The results are recorded here nevertheless.

#### Directly return by the wait strategy
| Implementation | ns/op|
|----------------|--|
|Busy Spin      |0.4013 ns/op|
|Yield          |0.4090 ns/op|
|Sleep wait     |0.4250 ns/op|

#### Actual timing for each implementation
| Implementation | ns/op|
|----------------|--|
|Busy Spin |0.2098 ns/op|
|Yield |43.79 ns/op|
|Sleep for 1 nano second |43.84 ns/op|

$-$

### Analysis
For the first set of benchmark result, the result is almost the same as there is no wait.
The time cost is mainly the comparison of the sequence, and the value is aligned with the benchmark result of sequence.

For the second set of benchmark, the execute of any statement is around 0.2 nanosecond, which means
the busy spin strategy will be wait/retry at 0.2 nanosecond interval.
Both Yield wait and Sleep wait cost around 43 nanosecond which is much larger.

$-$

## Sequencer
Sequencer is the backbone of the ring buffer, it maintains a cursor as a sequence pointing to the ring buffer,
and a list of gating sequences for each reader. By compare the values of the cursor and all gating sequence,
it coordinates the work of publisher and consumers.

$-$

### Benchmark results
Sequencer can claim 1 slot or many slots in a batch mode,

|Implementation| ns/op | 
|--|--|
|next() => claim 1 slot| 31.97 ns/op|
|nextN() => claim n slots| 3.109 ns/op|

$-$

### Analysis
Above test compares claiming 1 slot per call against claiming 10 slots per call
The batch mode  has a clear performance gain that the batch operation is cache friendly and fewer steps to be executed.

$-$

## RingBuffer



$-$

# Conclusion 



TODO: add some diagram/call graph etc ?
TODO: re-organize this MD 