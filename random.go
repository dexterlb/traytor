package traytor

import "math/rand"

// Random is a random generator with convenient methods
type Random struct {
	generator *rand.Rand
}

// NewRandom returns a new Random object initialized with the given seed
func NewRandom(seed int64) *Random {
	r := &Random{}
	source := rand.NewSource(seed)
	r.generator = rand.New(source)
	return r
}

// Float01 returns a random float between 0 and 1
func (r *Random) Float01() float64 {
	return r.generator.Float64()
}

// Float0A returns a random float between 0 and a
func (r *Random) Float0A(a float64) float64 {
	return r.Float01() * a
}

// FloatAB returns a random float between 0 and a
func (r *Random) FloatAB(a, b float64) float64 {
	return r.Float0A(b-a) + a
}
