package main

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// SimpleCircuit defines serial multiplication by constant

type Circuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y,public"`
}

// Define declares the circuit constraints
func (circuit *Circuit) Define(api frontend.API) error {
	product := api.Add(0, 1)
	for i := 0; i < 1048576; i++ {
		product = api.Mul(product, circuit.X)
	}
	api.AssertIsEqual(circuit.Y, product)
	return nil
}

func main() {
	// var p big.Int
	// p.SetString("28948022282369102195019208192891642640133654916059891165940443302932250623999", 10)
	// var f big.Int
	// f.SetString("3", 10)
	var ans big.Int
	ans.SetString("4428520108356670630000506301092116295361092437393728072026799124175523110981", 10) // take from gnark playground
	// for i := 0; i < 4; i++ {
	// 	ans = *ans.Mul(&ans, &f)
	// 	ans = *ans.Mod(&ans, &p)
	// }

	println("Y: ", ans.String())

	// compiles our circuit into a R1CS
	var circuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := Circuit{X: 3, Y: ans}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())

	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify

	// toggle on
	proofIci, _ := groth16.Prove(ccs, pk, witness, backend.WithIcicleAcceleration())
	groth16.Verify(proofIci, vk, publicWitness)
	// toggle off
	// proof, err := groth16.Prove(ccs, pk, witness)

	proof, _ := groth16.Prove(ccs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)
}
