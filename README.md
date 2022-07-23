# goDisruptor ![main status](https://github.com/chuchunf/goDisruptor/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/chuchunf/goDisruptor/branch/main/graph/badge.svg?token=XlzJA6ixJx)](https://codecov.io/gh/chuchunf/goDisruptor)

Another LMAX disruptor port in go 1.8 with generics support.

## Description
LMAX disruptor is a High Performance Inter-Thread Messaging Library as an alternative to bounded queueu, which makes use of padding to avoid memory false sharing, alignment of memory in stripe to be cache friendly etc. 

This is a port in Go with generics support, NOTE, Go's approach to concurrency is "**Don't communicate by sharing memory; share memory by communicating**". Channel is the prefered method for concurrency. This port follows disruptor's approach for better performance.

## Getting started

### Prerequisites
1. install latest go lang binary (1.18.3 and above) 
2. install latest vs code
 
### Running Testing
running all unit testing cases
```Shell
    go test -timeout 30s -tags -race goDisruptor/internal
    go test -timeout 30s -tags -race goDisruptor/pkg
    go test -timeout 30s -tags -race goDisruptor/example
```
running all performance testing benchmarks
```Shell
    go test -benchmem -run=^$ -tags -race -bench . goDisruptor/internal
    go test -benchmem -run=^$ -tags -race -bench . goDisruptor/pkg
    go test -benchmem -run=^$ -tags -race -bench . goDisruptor/example
```

### Usage
Please refer to [example](example) 

## Versioning
- 0.1
    - Initial Release

## Roadmap
- [x] support more waiting strategies 
- [ ] support publish N messages
- [ ] support consume N messages
- [ ] support mutiple consumers
- [ ] support proper DSL
- [ ] more examples

## License
Distributed under the MIT License. See [LICENSE.txt](LICENSE.txt) for more information.

## Acknowledgments
Thanks to the [LMAX-Disruptor](https://github.com/LMAX-Exchange/disruptor) project.