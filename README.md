# Benchmarks for ZK Data Center (ZaKi)


## Running benchmark

To run with ICICLE enabled, cd to project directory and

```sh
export CPATH=/usr/local/cuda/include
go mod tidy
go run -tags=icicle main.go
```

