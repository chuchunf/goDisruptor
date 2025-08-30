## WaitStrategy
WaitStrategy defines the strategy or mechanism for the producer or consumer to wait for the sequence.

Following waiting strategies are implemented and tested
* BusySpinWaitStrategy: waiting in an infinite loop 
* YieldWaitStrategy: yield the current thread's execution (via runtime.Gosched())
* SleepWaitStrategy: let current thread sleep for 1 Nano seconds (via time.Sleep())

#### Benchmark results
Obviously, the BusySpin strategy is the best performance implementation which has the best through put at a cost of wasting CPU cycle.
The results are recorded here nevertheless.

busysping 0.4013 ns/op
