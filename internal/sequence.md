## Sequence
Sequence is the fundemtal data strcture for the entire package.

In short, sequence is just a thread-safe counter, which for the consumer and producer to wait, get then process
a slot from the ring buffer. Refer to Sequencer for more details.

Following Sequences are tested
* struct with int64
* struct with [8]int65
* directly with int64

All performance testing are done in Windows 10 with i5-11400F @ 2.60GHz in console
```powershell
#bechmark
go test -benchmem -run=^$ -tags -race -bench ^BenchmarkSequenceGet$ goDisruptor/internal
#profile
go test -cpuprofile cpu.prof -memprofile mem.prof -benchmem -run=^$ -tags -race -bench ^BenchmarkSequenceGet$ goDisruptor/internal
#get top function calls
go tool pprof cpu.prof => top
# generate call graph
go tool pprof cpu.prof => png
# get top memory allocation
go tool pprof mem.prof => top
```

### Performance results 
||Get without GC|Get with GC|Set without GC|Set with GC|Concurrent Get and Set without GC
| -- | -- | -- | -- | -- | -- |
|struct int64|10.18 ns/op|11.42 ns/op|4.624 ns/op|4.615 ns/op|21.70 ns/op|
|struct [8]int64|23.10 ns/op|36.69 ns/op|4.714 ns/op|5.033 ns/op|49.09 ns/op|
|int64|0.2708 ns/op|0.2416 ns/op|4.936 ns/op|5.150 ns/op|5.291 ns/op|


### Get vs. Set

### Impact of GC
// slower without GC with vscoder, but faster when triggered directly
// no significant difference in call graph, likely due to the I/O, lock, Timer, scheduling etc.

### Impact of malloc
** mallocgc allocates 8 bytes cost around 8ns, compare to 20ns for 64 bytes
** atomic.StoreInt64 doesn't call mallocgc, it is consistent 4ns
** implementation using int64 directly, no allocation to heap, 0.2 ns only

### Impact of the struct
/*
** seems go lang allocate struct in heap, [8]int64 cost around 20ns
 */
get on direct int64 is much faster while set on direct int64 has no significant difference


### Impact of cache line / false sharing 

