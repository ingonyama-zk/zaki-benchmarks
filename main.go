package main

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const N = 20
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

func main() {
	// calculate Y
	var p big.Int
	p.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	var f big.Int
	f.SetString("42188824287", 10)
	var ans big.Int
	ans.SetString("1", 10)
	for i := 0; i < Size; i++ {
		ans = *ans.Mul(&ans, &f)
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
	assignment := Circuit{X: f, Y: ans}

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
}
