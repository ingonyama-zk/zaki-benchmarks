package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime/pprof"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const N = 24
const Size = 1 << N

type Circuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y,public"`
}

// Define declares the benchmark circuit: serial multiplication by constant
func (circuit *Circuit) Define(api frontend.API) error {
	product := api.Add(0, 1)
	for i := 0; i < Size; i++ {
		product = api.Mul(product, circuit.X)
	}
	api.AssertIsEqual(circuit.Y, product)
	return nil
}

var benchGPU bool
func init() {
	flag.BoolVar(&benchGPU, "bench_gpu", false, "Benchmarks GPU perfomance in addition to CPU performance")
}


func main() {
	flag.Parse()
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // Ensure the file is closed after the profiling is complete.

	// calculate Y
	var p big.Int
	p.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	var factor big.Int
	factor.SetString("42188824287", 10)
	var ans big.Int
	ans.SetString("1", 10)
	for i := 0; i < Size; i++ {
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
	// Start CPU profiling.
	if benchGPU {
		if err1 := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err1)
		}
		proofIci, err := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
		if err != nil {
			fmt.Println(err)
		}
		// Stop CPU profiling.
		pprof.StopCPUProfile()
		err = groth16.Verify(proofIci, vk, publicWitness)
		
		if err != nil {
			fmt.Println("Verify failed:", err)
		}

		proofIci2, err2 := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
		if err2 != nil {
			fmt.Println(err2)
		}
		err = groth16.Verify(proofIci2, vk, publicWitness)
		if err != nil {
			fmt.Println("Verify failed:", err)
		}
	}

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
