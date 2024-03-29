package maths

import (
	"encoding/json"
	"fmt"
	"math"
)

// Vec3 is a 3 dimensional vector
type Vec3 struct {
	X, Y, Z float64
}

// NewVec3 returns new 3 dimensional vector
func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{X: x, Y: y, Z: z}
}

// MakeZero makes all the dimentsions of the vector zero
func (v *Vec3) MakeZero() {
	v.X, v.Y, v.Z = 0, 0, 0
}

//NewVec3Array makes vector from array
func NewVec3Array(values [3]float64) *Vec3 {
	return &Vec3{
		X: values[0],
		Y: values[1],
		Z: values[2],
	}
}

//Length return the lenght of a vector
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// LengthSquared returns the square of the length of a vector
func (v *Vec3) LengthSquared() float64 {
	return (v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Scale multiplies all the dimension of the vector by the given multiplier
func (v *Vec3) Scale(multiplier float64) {
	v.X *= multiplier
	v.Y *= multiplier
	v.Z *= multiplier
}

// Add takes another vector and adds its dimensions to those of the given vector
func (v *Vec3) Add(other *Vec3) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

// Scaled returns a new Vec3 which is the product of the multiplication of the given vector and the multiplier
func (v *Vec3) Scaled(multiplier float64) *Vec3 {
	return NewVec3(
		v.X*multiplier,
		v.Y*multiplier,
		v.Z*multiplier,
	)
}

// Normalise sets the length to the given vector to 1
func (v *Vec3) Normalise() {
	v.SetLength(1.0)
}

// Normalised returns a new vector which length is 1
func (v *Vec3) Normalised() *Vec3 {
	normalisedVector := NewVec3(v.X, v.Y, v.Z)
	normalisedVector.SetLength(1.0)
	return normalisedVector
}

//Reflect makes the given vector equal to its reflected vector by the normal and is also normalised
func (v *Vec3) Reflect(normal *Vec3) {
	v.Normalise()
	v.Add(normal.Scaled(2 * DotProduct(normal, v.Negative())))
	v.Normalise()
}

//Reflected returns a new Vec3 witch is the normalised reflected vector of the given vector by the normal
func (v *Vec3) Reflected(normal *Vec3) *Vec3 {
	ray := v.Normalised()
	return AddVectors(ray, normal.Scaled(2*DotProduct(normal, ray.Negative()))).Normalised()
}

//Negative returns the opposite of the given vector
func (v *Vec3) Negative() *Vec3 {
	return v.Scaled(-1)
}

//GetDimension return X, Y, or Z depending on the given axis (or INF for wrong axis)
func (v *Vec3) GetDimension(axis int) float64 {
	switch axis {
	case Ox:
		return v.X
	case Oy:
		return v.Y
	case Oz:
		return v.Z
	}
	return Inf
}

//SetDimension makes the axis dimension of the vector equal to the value.
func (v *Vec3) SetDimension(axis int, value float64) {
	switch axis {
	case Ox:
		v.X = value
	case Oy:
		v.Y = value
	case Oz:
		v.Z = value
	}
}

//Negate makes the given vector equal to its opposite vector
func (v *Vec3) Negate() {
	v.Scale(-1)
}

//SetLength makes the lenght of the vector equal to the given newLength
func (v *Vec3) SetLength(newLength float64) {
	v.Scale(newLength / v.Length())
}

//String returns the string representation of the vector in the form of (x, y, z)
func (v *Vec3) String() string {
	return fmt.Sprintf("(%.3g, %.3g, %.3g)", SnapZero(v.X), SnapZero(v.Y), SnapZero(v.Z))
}

//AddVectors returns a new vector which is the sum of the two given vectors
func AddVectors(first, second *Vec3) *Vec3 {
	return NewVec3(first.X+second.X, first.Y+second.Y, first.Z+second.Z)
}

//MinusVectors returns a new vector which is the first vector minus the second
func MinusVectors(first, second *Vec3) *Vec3 {
	return NewVec3(first.X-second.X, first.Y-second.Y, first.Z-second.Z)
}

// DotProduct returns a float64 number which is the dot product of the two given vectors
func DotProduct(first, second *Vec3) float64 {
	return (first.X*second.X + first.Y*second.Y + first.Z*second.Z)
}

// CrossProduct returns a new Vec3 which is the cross product of the two given vectors
func CrossProduct(first, second *Vec3) *Vec3 {
	return NewVec3(
		first.Y*second.Z-first.Z*second.Y,
		first.Z*second.X-first.X*second.Z,
		first.X*second.Y-first.Y*second.X,
	)
}

//FaceForward flips the vector so that it's in the same hemispace as ray
func (v *Vec3) FaceForward(ray *Vec3) *Vec3 {
	if DotProduct(ray, v) < 0 {
		return v
	}
	return v.Negative()
}

//MixedProduct returns (a^b)*c
func MixedProduct(a, b, c *Vec3) float64 {
	return DotProduct(CrossProduct(a, b), c)
}

//Refract returns vector which is the refraction of incoming vector in relation with the normal with this ior
func Refract(incoming *Vec3, normal *Vec3, ior float64) *Vec3 {
	cosAlpha := DotProduct(incoming, normal)
	coeff := 1 - (ior*ior)*(1-cosAlpha*cosAlpha)
	if coeff < 0 { // total inner reflection - refracted vector is undefined
		return nil
	}
	return MinusVectors(
		incoming.Scaled(ior),
		normal.Scaled((ior*cosAlpha + math.Sqrt(coeff))),
	)
}

//UnmarshalJSON implements the json.Unmarshaler interface
func (v *Vec3) UnmarshalJSON(data []byte) error {
	var unmarshaled [3]float64
	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}
	v.X = unmarshaled[0]
	v.Y = unmarshaled[1]
	v.Z = unmarshaled[2]
	return nil
}
