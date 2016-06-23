package mesh

import (
	"fmt"
	"testing"

	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
	"github.com/stretchr/testify/assert"
)

func ExampleBoundingBox_AddPoint() {
	box := NewBoundingBox()

	box.AddPoint(maths.NewVec3(-1, -1, -1))
	box.AddPoint(maths.NewVec3(0, 5, -0.5))
	box.AddPoint(maths.NewVec3(1, 0, 2))

	fmt.Printf("%s\n", box)

	// Output:
	// bbox[min: (-1, -1, -1), max: (1, 5, 2)]
}

func ExampleBoundingBox_Inside() {
	box := &BoundingBox{
		MinVolume: [3]float64{-1, -1, -1},
		MaxVolume: [3]float64{1, 1, 1},
	}

	fmt.Printf("0, 0, 0: %v\n", box.Inside(maths.NewVec3(0, 0, 0)))
	fmt.Printf("0, 0, 2: %v\n", box.Inside(maths.NewVec3(0, 0, 2)))

	// Output:
	// 0, 0, 0: true
	// 0, 0, 2: false
}

func ExampleBoundingBox_Intersect() {
	box := &BoundingBox{
		MinVolume: [3]float64{-1, -1, -1},
		MaxVolume: [3]float64{1, 1, 1},
	}

	ray1 := &ray.Ray{
		Start:     *maths.NewVec3(0, 0, 0),
		Direction: *maths.NewVec3(1, 0, 0),
	}

	ray2 := &ray.Ray{
		Start:     *maths.NewVec3(-5, 0, 0.5),
		Direction: *maths.NewVec3(1, 0, 0),
	}

	ray3 := &ray.Ray{
		Start:     *maths.NewVec3(-5, 0, 0.5),
		Direction: *maths.NewVec3(-1, 0, 0),
	}

	ray4 := &ray.Ray{
		Start:     *maths.NewVec3(-5, 0, 0),
		Direction: *maths.NewVec3(1, 0, 5),
	}
	ray1.Init()
	ray2.Init()
	ray3.Init()
	ray4.Init()

	fmt.Printf("ray 1: %v\n", box.Intersect(ray1))
	fmt.Printf("ray 2: %v\n", box.Intersect(ray2))
	fmt.Printf("ray 3: %v\n", box.Intersect(ray3))
	fmt.Printf("ray 4: %v\n", box.Intersect(ray4))

	// Output:
	// ray 1: true
	// ray 2: true
	// ray 3: false
	// ray 4: false
}

func TestIntersectBoundingBoxZeroVolume(t *testing.T) {
	box := NewBoundingBox()
	box.AddPoint(maths.NewVec3(0, 0, 0))
	box.AddPoint(maths.NewVec3(1, 0, 0))
	box.AddPoint(maths.NewVec3(0, 1, 0))
	box.AddPoint(maths.NewVec3(1, 1, 0))

	ray := &ray.Ray{
		Start:     *maths.NewVec3(0.5, 0.5, 1),
		Direction: *maths.NewVec3(0, 0, -1),
	}

	if !box.Intersect(ray) {
		t.Error("ray.Ray must intersect box")
	}
}

func TestBoundingBoxSplit(t *testing.T) {
	box := &BoundingBox{
		MinVolume: [3]float64{-3.18, -3.96, -1.77},
		MaxVolume: [3]float64{4.58, 1.74, 4.56},
	}

	ray := &ray.Ray{
		Start:     *maths.NewVec3(4.49, -7.3, 5.53),
		Direction: *maths.NewVec3(-0.60, 0.5, -0.62),
	}
	ray.Init()

	if !box.Intersect(ray) {
		t.Error("ray should intersect big box")
	}

	left, right := box.Split(2, 1.39)

	assertEqualVectors(t, maths.NewVec3(-3.18, -3.96, -1.77), maths.NewVec3Array(left.MinVolume))
	assertEqualVectors(t, maths.NewVec3(-3.18, -3.96, 1.39), maths.NewVec3Array(right.MinVolume))
	assertEqualVectors(t, maths.NewVec3(4.58, 1.74, 1.39), maths.NewVec3Array(left.MaxVolume))
	assertEqualVectors(t, maths.NewVec3(4.58, 1.74, 4.56), maths.NewVec3Array(right.MaxVolume))

	if box.IntersectWall(2, 1.39, ray) {
		t.Error("ray shouldn't intersect the wall")
	}

	if right.Intersect(ray) {
		t.Error("ray shouldn't intersect right")
	}

	if !left.Intersect(ray) {
		t.Error("ray should intersect left")
	}
}

func assertEqualVectors(t *testing.T, expected *maths.Vec3, v *maths.Vec3) {
	assert := assert.New(t)
	assert.InDelta(expected.X, v.X, maths.Epsilon)
	assert.InDelta(expected.Y, v.Y, maths.Epsilon)
	assert.InDelta(expected.Z, v.Z, maths.Epsilon)
}
