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

//FloatToInt return an int between 0 and 1 constructed from a given float between 0 and 255
func FloatToInt(x float64) int {
	a := 0.055
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 255
	}
	if x <= 0.00313008 {
		x = x * 12.02
	} else {
		x = (1.0+a)*math.Pow(x, 1.0/2.4) - a
	}
	return Round(x * 255.0)
}
