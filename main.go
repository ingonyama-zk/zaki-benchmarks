package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type Circuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y,public"`
}

// Define declares the benchmark circuit: serial multiplication by constant
func (circuit *Circuit) Define(api frontend.API) error {
	product := api.Add(0, 1)
	size := 1 << N
	for i := 0; i < size; i++ {
		product = api.Mul(product, circuit.X)
	}
	api.AssertIsEqual(circuit.Y, product)
	return nil
}

var benchGPU, benchCPU, benchAll bool
var N int
func init() {
	flag.BoolVar(&benchGPU, "bench_gpu", false, "Benchmarks GPU perfomance. Default: false")
	flag.BoolVar(&benchCPU, "bench_cpu", false, "Benchmarks CPU perfomance. Default: false")
	flag.BoolVar(&benchAll, "bench_all", false, "Benchmarks GPU and CPU perfomance. Default: false")
	flag.IntVar(&N, "size", 24, "Size as a power of two that should be benched; e.g. 20 for benching 2^20.")
}

func main() {
	flag.Parse()

	if !benchCPU && !benchGPU && !benchAll {
		panic("No hardware was selected for benchmarking. Please select CPU, GPU, or both with the flags --bench_gpu, --bench_cpu, or --bench_all")
	}

	// calculate Y
	var p big.Int
	p.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	var factor big.Int
	factor.SetString("42188824287", 10)
	var ans big.Int
	ans.SetString("1", 10)
	size := 1 << N
	for i := 0; i < size; i++ {
		ans = *ans.Mul(&ans, &factor)
		ans = *ans.Mod(&ans, &p)
	}

	// compiles the circuit into a R1CS
	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println(err)
	}

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println(err)
	}

	// witness definition
	assignment := Circuit{X: factor, Y: ans}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		fmt.Println(err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		fmt.Println(err)
	}

	// Prove & Verify

	// on GPU
	if benchGPU || benchAll {
		proofIci, err := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
		if err != nil {
			fmt.Println(err)
		}
		err = groth16.Verify(proofIci, vk, publicWitness)

		if err != nil {
			fmt.Println("Verify failed:", err)
		}

		os.Setenv("profile", "ON")
		proofIci2, err2 := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
		if err2 != nil {
			fmt.Println(err2)
		}
		os.Unsetenv("profile")
		err = groth16.Verify(proofIci2, vk, publicWitness)
		if err != nil {
			fmt.Println("Verify failed:", err)
		}
	}

	if benchCPU || benchAll {
		// on CPU
		proof, err := groth16.Prove(ccs, pk, witness)
		if err != nil {
			fmt.Println(err)
		}
		err = groth16.Verify(proof, vk, publicWitness)
		if err != nil {
			fmt.Println("Verify failed:", err)
		}
	}
}
