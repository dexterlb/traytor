package random

import (
	"math"
	"math/rand"

	"github.com/DexterLB/traytor/maths"
)

// Random is a random generator with convenient methods
type Random struct {
	generator *rand.Rand
}

// New returns a new Random object initialized with the given seed
func New(seed int64) *Random {
	r := &Random{}
	source := rand.NewSource(seed)
	r.generator = rand.New(source)
	return r
}

// Vec3Sphere returns a random unit vector
func (r *Random) Vec3Sphere() *maths.Vec3 {
	u := r.FloatAB(-1, 1)
	theta := r.Float02Pi()
	return maths.NewVec3(
		math.Sqrt(1-u*u)*math.Cos(theta),
		math.Sqrt(1-u*u)*math.Sin(theta),
		u,
	)
}

// Vec3Hemi returns a random unit vector in the hemisphere defined by normal
func (r *Random) Vec3Hemi(normal *maths.Vec3) *maths.Vec3 {
	vec := r.Vec3Sphere()
	if maths.DotProduct(vec, normal) < 0 {
		vec.Negate()
	}
	return vec
}

// Vec3HemiCos returns a random unit vector chosen on a
// cosine-weighed hemisphere defined by normal
func (r *Random) Vec3HemiCos(normal *maths.Vec3) *maths.Vec3 {
	ox := maths.CrossProduct(maths.NewVec3(42, 56, -15), normal)
	for math.Abs(ox.Length()) < maths.Epsilon {
		ox = maths.CrossProduct(r.Vec3Sphere(), normal)
	}

	oy := maths.CrossProduct(ox, normal)
	ox.Normalise()
	oy.Normalise()

	u := r.Float01()
	radius := math.Sqrt(u)
	theta := r.Float02Pi()

	vec := normal.Scaled(math.Sqrt(math.Max(0, 1-u)))
	vec.Add(ox.Scaled(radius * math.Cos(theta)))
	vec.Add(oy.Scaled(radius * math.Sin(theta)))

	return vec
}

// Float01 returns a random float between 0 and 1
func (r *Random) Float01() float64 {
	return r.generator.Float64()
}

// Float0Pi returns a random float between 0 and Pi
func (r *Random) Float0Pi() float64 {
	return r.generator.Float64() * math.Pi
}

// Float02Pi returns a random float between 0 and 2*Pi
func (r *Random) Float02Pi() float64 {
	return r.generator.Float64() * 2 * math.Pi
}

// Float0A returns a random float between 0 and a
func (r *Random) Float0A(a float64) float64 {
	return r.Float01() * a
}

// FloatAB returns a random float between 0 and a
func (r *Random) FloatAB(a, b float64) float64 {
	return r.Float0A(b-a) + a
}

// Int640N returns a random int64 within [0..n]
func (r *Random) Int640N(n int64) int64 {
	return r.generator.Int63n(n + 1)
}

// Int320N returns a random int32 within [0..n]
func (r *Random) Int320N(n int32) int32 {
	return r.generator.Int31n(n + 1)
}

// Int0N returns a random int within [0..n]
func (r *Random) Int0N(n int) int {
	return r.generator.Intn(n + 1)
}

// Int64AB returns a random int64 within [a..b]
func (r *Random) Int64AB(a, b int64) int64 {
	return r.Int640N(b-a) + a
}

// Int32AB returns a random int32 within [a..b]
func (r *Random) Int32AB(a, b int32) int32 {
	return r.Int320N(b-a) + a
}

// IntAB returns a random int within [a..b]
func (r *Random) IntAB(a, b int) int {
	return r.Int0N(b-a) + a
}

// Bool returns true or false at random
func (r *Random) Bool() bool {
	return (r.Int0N(1) == 0)
}

// Sign returns -1 or 1 at random
func (r *Random) Sign() int {
	if r.Bool() {
		return 1
	}
	return -1
}

// Sign32 returns -1 or 1 at random
func (r *Random) Sign32() int32 {
	if r.Bool() {
		return 1
	}
	return -1
}

// Sign64 returns -1 or 1 at random
func (r *Random) Sign64() int64 {
	if r.Bool() {
		return 1
	}
	return -1
}

// NewSeed returns a random seed
func (r *Random) NewSeed() int64 {
	return r.generator.Int63()
}
