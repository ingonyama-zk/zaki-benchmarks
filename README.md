# Benchmarks for ZK Data Center (ZaKi)


## Running benchmark

To run with ICICLE enabled, cd to project directory and

```sh
go run -tags=icicle main.go
```

If you don't need GPU (CPU) comment out section `on GPU` (`on CPU`) in `main.go`:


```go
    // on GPU
	proofIci, err := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
	if err != nil {
		fmt.Println(err)
	}
	groth16.Verify(proofIci, vk, publicWitness)

	// on CPU
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Println(err)
	}
	groth16.Verify(proof, vk, publicWitness)
```

