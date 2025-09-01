## Sequencer
Sequencer is the backbone of the ring buffer, it maintains a cursor as a sequence pointing to the ring buffer, 
and a list of gating sequences for each reader. By compare the values of the cursor and all gating sequence, 
it coordinates the work of publisher and consumers.

### Benchmark results
Sequencer can claim 1 slot or many slots in a batch mode, 

|Implementation| ns/op | 
|--|--|
|next() => claim 1 slot| 31.97 ns/op|
|nextN() => claim n slots| 3.109 ns/op|

### Analysis
Above test compares claiming 1 slot per call against claiming 10 slots per call
The batch mode  has a clear performance gain that the batch operation is cache friendly and fewer steps to be executed.