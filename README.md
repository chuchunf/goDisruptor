# goDisruptor ![main status](https://github.com/chuchunf/goDisruptor/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/chuchunf/goDisruptor/branch/main/graph/badge.svg?token=XlzJA6ixJx)](https://codecov.io/gh/chuchunf/goDisruptor)

Another LMAX disruptor port in go 1.8 with generics support.

## TODO: 
* mutiple writer support
* diamond setup/example
* find hot path 
    * waitFor use reference not value
    * use primitive int for count instead of [8]int64 
* readme update
* proper write up 