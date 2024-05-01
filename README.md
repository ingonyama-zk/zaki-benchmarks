# Benchmarks for ZK Data Center (ZaKi)


## Running benchmark

To run with ICICLE enabled, cd to project directory and

```sh
export CPATH=/usr/local/cuda/include
go mod tidy
go run -tags=icicle main.go <options>
```

### USAGE
```sh
go run -tags=icicle main.go --help

  -bench_all
        Benchmarks GPU and CPU perfomance. Default: false
  -bench_cpu
        Benchmarks CPU perfomance. Default: false
  -bench_gpu
        Benchmarks GPU perfomance. Default: false
  -size int
        Size as a power of two that should be benched; e.g. 20 for benching 2^20 (default 24)
  -profile
        Prints profile timings. Default: false
```
