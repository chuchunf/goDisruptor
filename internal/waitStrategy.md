## WaitStrategy
WaitStrategy defines the strategy or mechanism for the producer or consumer to wait for the sequence.

Following waiting strategies are implemented and tested
* BusySpinWaitStrategy: waiting in an infinite loop 
* YieldWaitStrategy: yield the current thread's execution (via runtime.Gosched())
* SleepWaitStrategy: let current thread sleep for 1 Nano seconds (via time.Sleep())

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


### Analysis
For the first set of benchmark result, the result is almost the same as there is no wait.
The time cost is mainly the comparison of the sequence, and the value is aligned with the benchmark result of sequence.

For the second set of benchmark, the execute of any statement is around 0.2 nanosecond, which means
the busy spin strategy will be wait/retry at 0.2 nanosecond interval.
Both Yield wait and Sleep wait cost around 43 nanosecond which is much larger.

