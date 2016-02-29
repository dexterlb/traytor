package traytor

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectors(t *testing.T) {
	assert := assert.New(t)

	rnd := NewRandom(42)

	vec := rnd.Vec3Sphere()

	assert.InDelta(1, vec.Length(), Epsilon, "Random unit vector's length should be 1")

	norm := vec

	vec = rnd.Vec3Hemi(norm)
	if DotProduct(vec, norm) < 0 {
		t.Error("Random hemi vector is in the wrong hemisphere")
	}

	vec = rnd.Vec3HemiCos(norm)
	if DotProduct(vec, norm) < 0 {
		t.Error("Random cosine-weighed hemi vector is in the wrong hemisphere")
	}
}

func TestFloats(t *testing.T) {
	rnd := NewRandom(42)
	number := rnd.Float01()
	if number < 0 || number > 1 {
		t.Error("Float01 should be in [0..1]")
	}

	number = rnd.Float0A(15.0)
	if number < 0 || number > 15.0 {
		t.Error("Float0A should be in [0..a]")
	}

	number = rnd.FloatAB(-10.0, 5.0)
	if number < -10.0 || number > 5.0 {
		t.Error("FloatAB should be in [a..b]")
	}

	number = rnd.Float0Pi()
	if number < 0 || number > math.Pi {
		t.Error("Float0Pi should be in [0..pi]")
	}

	number = rnd.Float02Pi()
	if number < 0 || number > 2*math.Pi {
		t.Error("Float02Pi should be in [0..2pi]")
	}
}

func TestInts(t *testing.T) {
	rnd := NewRandom(42)

	number32 := rnd.Int320N(15)
	if number32 < 0 || number32 > 15 {
		t.Error("Int320N should be in [0..n]")
	}

	number32 = rnd.Int32AB(-10, 5)
	if number32 < -10 || number32 > 5 {
		t.Error("Int32AB should be in [a..b]")
	}

	number64 := rnd.Int640N(15)
	if number64 < 0 || number64 > 15 {
		t.Error("Int640N should be in [0..n]")
	}

	number64 = rnd.Int64AB(-10, 5)
	if number64 < -10 || number64 > 5 {
		t.Error("Int64AB should be in [a..b]")
	}

	number := rnd.Int0N(15)
	if number < 0 || number > 15 {
		t.Error("Int0N should be in [0..n]")
	}

	number = rnd.IntAB(-10, 5)
	if number < -10 || number > 5 {
		t.Error("IntAB should be in [a..b]")
	}
}

func TestBools(t *testing.T) {
	rnd := NewRandom(42)

	b := rnd.Bool()
	if b != true && b != false {
		t.Error("The world is going to end.")
	}

	sign := rnd.Sign()
	if sign != -1 && sign != 1 {
		t.Error("Sign should be -1 or 1")
	}

	sign32 := rnd.Sign32()
	if sign32 != -1 && sign32 != 1 {
		t.Error("Sign should be -1 or 1")
	}

	sign64 := rnd.Sign64()
	if sign64 != -1 && sign64 != 1 {
		t.Error("Sign should be -1 or 1")
	}
}

func TestSeed(t *testing.T) {
	rnd1 := NewRandom(42)
	rnd2 := NewRandom(42)
	rnd3 := NewRandom(56)
	n1 := rnd1.Float01()
	n2 := rnd2.Float01()
	n3 := rnd3.Float01()
	if n1 != n2 {
		t.Error("First values obtained with the same seed must be the same")
	}
	if n2 == n3 {
		t.Error("First values obtained with different seeds are very unlikely to be the same")
	}
	n12 := rnd1.Float01()
	n22 := rnd2.Float01()
	if n12 != n22 {
		t.Error("Second values obtained with the same seed must be the same")
	}
}
