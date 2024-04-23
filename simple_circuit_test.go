package main

import (
	"testing"
	"github.com/consensys/gnark/test"
)

func TestCubicEquation(t *testing.T) {
	assert := test.NewAssert(t)

	var simpleCircuit SimpleCircuit

	assert.ProverFailed(&simpleCircuit, &SimpleCircuit{
		X: 42,
		Y: 42,
	})

	assert.ProverSucceeded(&simpleCircuit, &SimpleCircuit{
		X: 3,
		Y: 35,
	})
}