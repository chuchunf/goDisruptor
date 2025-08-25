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
| int64           | 0.2119 ns/op   | 0.4066 ns/op | 1.612 ns/op   | 1.611 ns/op   |


### Impact of GC
with GC turns off, it

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
