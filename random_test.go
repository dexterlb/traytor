package traytor

import "testing"

func TestFloats(t *testing.T) {
	rnd := NewRandom(42)
	var number float64 = rnd.Float01()
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
}

func TestInts(t *testing.T) {
	rnd := NewRandom(42)

	var number32 int32 = rnd.Int320N(15)
	if number32 < 0 || number32 > 15 {
		t.Error("Int320N should be in [0..n]")
	}

	number32 = rnd.Int32AB(-10, 5)
	if number32 < -10 || number32 > 5 {
		t.Error("Int32AB should be in [a..b]")
	}

	var number64 int64 = rnd.Int640N(15)
	if number64 < 0 || number64 > 15 {
		t.Error("Int640N should be in [0..n]")
	}

	number64 = rnd.Int64AB(-10, 5)
	if number64 < -10 || number64 > 5 {
		t.Error("Int64AB should be in [a..b]")
	}

	var number int = rnd.Int0N(15)
	if number < 0 || number > 15 {
		t.Error("Int0N should be in [0..n]")
	}

	number = rnd.IntAB(-10, 5)
	if number < -10 || number > 5 {
		t.Error("IntAB should be in [a..b]")
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
