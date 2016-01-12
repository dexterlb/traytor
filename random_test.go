package traytor

import "testing"

func TestFloat01(t *testing.T) {
	rnd := NewRandom(42)
	number := rnd.Float01()
	if number < 0 || number > 1 {
		t.Error("Float01 should be in [0..1]")
	}
}

func TestFloat0A(t *testing.T) {
	rnd := NewRandom(42)
	number := rnd.Float0A(15.0)
	if number < 0 || number > 15.0 {
		t.Error("Float0A should be in [0..a]")
	}
}

func TestFloatAB(t *testing.T) {
	rnd := NewRandom(42)
	number := rnd.FloatAB(-10.0, 5.0)
	if number < -10.0 || number > 5.0 {
		t.Error("FloatAB should be in [a..b]")
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
