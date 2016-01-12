package traytor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleVec3_ToZero() {
	v := NewVec3(1, 2, 3)
	fmt.Printf("%s\n", v)
	v.ToZero()
	fmt.Printf("%s\n", v)

	// Output:
	// (1, 2, 3)
	// (0, 0, 0)
	//
}

func ExampleVec3_Length() {
	v := NewVec3(1, 2, 2)
	fmt.Printf("%.3g\n", v.Length())

	// Output:
	// 3
	//
}

func ExampleVec3_LengthSquared() {
	v := NewVec3(1, 2, 2)
	fmt.Printf("%.3g\n", v.LengthSquared())

	// Output:
	// 9
	//
}

func ExampleVec3_Scale() {
	v := NewVec3(1, -2, 3)
	v.Scale(2)
	fmt.Printf("%s\n", v)

	// Output:
	// (2, -4, 6)
	//
}

func ExampleVec3_Add() {
	v := NewVec3(1, -2, 3)
	fmt.Printf("%s\n", v)
	v.Add(NewVec3(1, 2, 3))
	fmt.Printf("%s\n", v)

	// Output:
	// (1, -2, 3)
	// (2, 0, 6)
	//
}

func asserEqualVectors(t *testing.T, expected *Vec3, v *Vec3) {
	assert := assert.New(t)
	assert.InEpsilon(expected.X, v.X, Epsilon)
	assert.InEpsilon(expected.Y, v.Y, Epsilon)
	assert.InEpsilon(expected.Z, v.Z, Epsilon)
}

func TestScaling(t *testing.T) {
	v := NewVec3(1, 0, 3)
	scaled := v.Scaled(2)

	asserEqualVectors(t, NewVec3(2, 0, 6), scaled)
}

func TestNormalising(t *testing.T) {
	assert := assert.New(t)
	v := NewVec3(1, 2, 3)
	normalised := v.Normalised()
	v.Normalise()

	assert.InEpsilon(v.Length(), 1, Epsilon, "Normalising should make vector's lenght 1")
	assert.InEpsilon(1, normalised.Length(), Epsilon, "Normalised should return vector with length 1")
}

func TestReflecting(t *testing.T) {
	normal := NewVec3(0, 1, 0)
	v := NewVec3(1, -1, 0)
	reflected := v.Reflected(normal)
	v.Reflect(normal)

	asserEqualVectors(t, NewVec3(1, 1, 0).Normalised(), v)
	asserEqualVectors(t, NewVec3(1, 1, 0).Normalised(), reflected)
}
