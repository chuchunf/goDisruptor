# goDisruptor ![main status](https://github.com/chuchunf/goDisruptor/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/chuchunf/goDisruptor/branch/main/graph/badge.svg?token=XlzJA6ixJx)](https://codecov.io/gh/chuchunf/goDisruptor)

Another LMAX disruptor port in go 1.8 with generics support.

## TODO: 
* why stoppging GC has different effect ?
* find out why PIN cpu is slow
* cache memory alignment (cachegrind)
* false sharing
* find hot path 
* publishN example
* mutiple writer support
* diamond setup/example
* readme update
* proper write up 