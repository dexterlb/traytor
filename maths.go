package traytor

import "math"

const (
	// Epsilon is a very small number
	Epsilon float64 = 1e-9
	Inf     float64 = 1e99
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
