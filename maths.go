package traytor

import "math"

const (
	// Epsilon is a very small number
	Epsilon float64 = 1e-9
)

//Round returns the nearest int to a given float number
func Round(number float64) int {
	return int(math.Floor(number + 0.5))
}
