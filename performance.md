
## Introduction
The LMAX Disruptor is a high-performance inter-thread messaging library developed in Java, renowned for its exceptionally low latency  
by using lock-free algorithms, careful memory layout and other consideration at hardware level to minimize contention and improve the throughput.

The objective of the port is to learn to implement, test and verify low latency library in general, as golang prefers goroutine.
As Go provides powerful built-in primitives like channels and goroutines.

This document focuses on the performance analysis of this Go port of the LMAX Disruptor. 

$-$

## Design of LMAX-Disruptor
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

## Benchmark and Performance Analysis

### Sequence


### WaitStrategy


### SequenceBarrier


### Sequencer


### RingBuffer



$-$

## Conclusion 
















#################################
