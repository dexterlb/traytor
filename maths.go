package traytor

import "math"

const (
	// Ox represents the X axis
	Ox = iota
	// Oy represents the Y axis
	Oy
	// Oz represents the Z axis
	Oz
	// Leaf represents a leaf in a tree enumerated with axes
	Leaf
)

const (
	// Epsilon is a very small number
	Epsilon float64 = 1e-9
	// Inf is a very large number
	Inf float64 = 1e99
)

// Round returns the nearest int to a given float number
func Round(number float64) int {
	return int(math.Floor(number + 0.5))
}

// Round32 returns the nearest int to a given float32 number
func Round32(number float32) int {
	return int(math.Floor(float64(number + 0.5)))
}

// Pow32 is a Pow function which uses float32
func Pow32(x, a float32) float32 {
	return float32(math.Pow(float64(x), float64(a)))
}

// SolveEquation solves the following system (returning x and y):
// | a1 * x + b1 * y + c1 = 0
// | a2 * x + b2 * y + c2 = 0
func SolveEquation(a, b, c *Vec3) (float64, float64) {
	coefficientMatrix := [2][2]float64{{a.X, b.X}, {a.Y, b.Y}}
	constantCoefficientMatrix := [2]float64{c.X, c.Y}

	det := coefficientMatrix[0][0]*coefficientMatrix[1][1] -
		coefficientMatrix[1][0]*coefficientMatrix[0][1]
	x := (constantCoefficientMatrix[0]*coefficientMatrix[1][1] -
		constantCoefficientMatrix[1]*coefficientMatrix[0][1]) / det
	y := (constantCoefficientMatrix[1]*coefficientMatrix[0][0] -
		constantCoefficientMatrix[0]*coefficientMatrix[1][0]) / det
	return x, y
}

// SnapZero returns 0 if x is too close to 0, and returns x otherwise
func SnapZero(x float64) float64 {
	if x < Epsilon && x > -Epsilon {
		return 0
	}
	return x
}

//Between check whether point is between min and max (+- Epsilon)
func Between(min, max, point float64) bool {
	return (min-Epsilon <= point && max+Epsilon >= point)
}
