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

### Performance results 
|                 | Get without GC | Get with GC  | Set without GC | Set with GC   | 
|-----------------|----------------|--------------|---------------|---------------|
| struct int64    | 0.2182 ns/op   | 0.2119 ns/op | 1.610 ns/op   | 1.688 ns/op   | 
| struct [8]int64 | 0.2142 ns/op   | 0.2194 ns/op | 1.586 ns/op   | 1.607 ns/op   |
| int64           | 0.2119 ns/op   | 0.2121 ns/op | 1.612 ns/op   | 1.611 ns/op   |


### Impact of sampling rate
Given we're testing for nano second changes and the function used to run less than 1 second,
the default sample rate of 100 is too small. 
For all the benchmarking, we will use sample rate at 10000
```go
 runtime.SetCPUProfileRate(10000)
```

### Impact of GC
From the benchmark above, the impact of GC in all the cases are negligible, probably due to the fact that 
the memory footprint and change are relatively small, there are not much GC activities anyway.

This can be confirmed by the cpu.prof that the call graph are almost identical.
Although it needs to be noted that from the mem.prof, the tests without GC allocate a larger memory initially.

### Get vs. Set

### Impact of malloc
** malloc allocates 8 bytes cost around 8ns, compare to 20ns for 64 bytes
** atomic.StoreInt64 doesn't call malloc, it is consistent 4ns
** implementation using int64 directly, no allocation to heap, 0.2 ns only

### Impact of the struct
/*
** seems go lang allocate struct in heap, [8]int64 cost around 20ns
 */
get on direct int64 is much faster while set on direct int64 has no significant difference


### Impact of cache line / false sharing 
how to confirm this false sharing ?
write a testing that have 2 sequences next to each other and use to get/set


TODO: use testprofile instead call the function SetCPUProfileRate ?