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

	assert.InEpsilon(1, v.Length(), Epsilon, "Normalising should make vector's lenght 1")
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

func TestNegative(t *testing.T) {
	v := NewVec3(1, -2, 0)
	negative := v.Negative()
	v.Negate()

	asserEqualVectors(t, NewVec3(-1, 2, 0), v)
	asserEqualVectors(t, NewVec3(-1, 2, 0), negative)
}

func ExampleVec3_Negate() {
	v := NewVec3(1, -2, 10)
	fmt.Printf("%s\n", v)
	v.Negate()
	fmt.Printf("%s\n", v)

	// Output:
	// (1, -2, 10)
	// (-1, 2, -10)
	//
}

func TestSettingLength(t *testing.T) {
	assert := assert.New(t)
	v := NewVec3(3, 4, 0)
	assert.InEpsilon(5, v.Length(), Epsilon, "Vector's length should be 5")
	v.SetLength(10)
	asserEqualVectors(t, NewVec3(6, 8, 0), v)
}

func TestAddingVectors(t *testing.T) {
	sum := AddVectors(NewVec3(1, 2, 3), NewVec3(2, -2, 3))
	asserEqualVectors(t, NewVec3(3, 0, 6), sum)
}

func ExampleAddVectors() {
	fmt.Printf("v = AddVectors((1, 2, 3), (2, -2, 3))\n")
	v := AddVectors(NewVec3(1, 2, 3), NewVec3(2, -2, 3))
	fmt.Printf("v = %s\n", v)

	// Output:
	// v = AddVectors((1, 2, 3), (2, -2, 3))
	// v = (3, 0, 6)
	//
}

func TestMinusingVectors(t *testing.T) {
	sum := MinusVectors(NewVec3(1, 2, 3), NewVec3(2, 2, 2))
	asserEqualVectors(t, NewVec3(-1, 0, 1), sum)
}

func ExampleMinusVectors() {
	fmt.Printf("v = MinusVectors((1, 2, 3), (2, 2, 2))\n")
	v := MinusVectors(NewVec3(1, 2, 3), NewVec3(2, 2, 2))
	fmt.Printf("v = %s\n", v)

	// Output:
	// v = MinusVectors((1, 2, 3), (2, 2, 2))
	// v = (-1, 0, 1)
	//
}

func TestDotProduct(t *testing.T) {
	assert := assert.New(t)
	dotProduct := DotProduct(NewVec3(1, 2, 1), NewVec3(-1, 0, 43))
	assert.InEpsilon(42, dotProduct, Epsilon)
}

func ExampleDotProduct() {
	fmt.Printf("p = DotProduct((1, 2, 1), (-1, 0, 43))\n")
	p := DotProduct(NewVec3(1, 2, 1), NewVec3(-1, 0, 43))
	fmt.Printf("p = %.3g\n", p)

	// Output:
	// p = DotProduct((1, 2, 1), (-1, 0, 43))
	// p = 42
	//
}

func TestCrossProduct(t *testing.T) {
	crossProduct := CrossProduct(NewVec3(1, 0, 0), NewVec3(0, 1, 0))
	asserEqualVectors(t, NewVec3(0, 0, 1), crossProduct)
}

func ExampleCrossProduct() {
	fmt.Printf("p = CrossProduct((1, 0, 0), (0, 1, 0))\n")
	p := CrossProduct(NewVec3(1, 0, 0), NewVec3(0, 1, 0))
	fmt.Printf("p = %s\n", p)

	// Output:
	// p = CrossProduct((1, 0, 0), (0, 1, 0))
	// p = (0, 0, 1)
	//
}

func TestFacingForward(t *testing.T) {
	normal := NewVec3(0, 0, -1)
	newNormal := normal.FaceForward(NewVec3(0, 0, 1))
	asserEqualVectors(t, NewVec3(0, 0, -1), newNormal)
	normal = NewVec3(0, 0, 1)
	newNormal = normal.FaceForward(NewVec3(0, 0, 1))
	asserEqualVectors(t, NewVec3(0, 0, -1), newNormal)

}
