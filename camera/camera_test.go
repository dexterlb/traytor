package camera

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DexterLB/traytor/maths"
	"github.com/stretchr/testify/assert"
)

func ExamplePinholeCamera_ShootRay() {
	var c Camera = &PinholeCamera{
		Focus:      *maths.NewVec3(0, 0, 0),
		TopLeft:    *maths.NewVec3(-1, 1, 1),
		TopRight:   *maths.NewVec3(1, 1, 1),
		BottomLeft: *maths.NewVec3(-1, 1, -1),
	}

	ray := c.ShootRay(0, 0)
	fmt.Printf("0, 0: %s\n", ray)

	ray = c.ShootRay(0.5, 0.5)
	fmt.Printf("0.5, 0.5: %s\n", ray)

	// Output:
	// 0, 0: (0, 0, 0) -> (-0.577, 0.577, 0.577)
	// 0.5, 0.5: (0, 0, 0) -> (0, 1, 0)
}

func ExamplePinholeCamera_ShootRay_weirdfocus() {
	var c Camera = &PinholeCamera{
		Focus:      *maths.NewVec3(-2, 15, 3),
		TopLeft:    *maths.NewVec3(-3, 16, 4),
		TopRight:   *maths.NewVec3(-1, 16, 4),
		BottomLeft: *maths.NewVec3(-3, 16, 2),
	}

	ray := c.ShootRay(0, 0)
	fmt.Printf("0, 0: %s\n", ray)

	ray = c.ShootRay(0.5, 0.5)
	fmt.Printf("0.5, 0.5: %s\n", ray)

	// Output:
	// 0, 0: (-2, 15, 3) -> (-0.577, 0.577, 0.577)
	// 0.5, 0.5: (-2, 15, 3) -> (0, 1, 0)
}

func TestPinholeCameraJson(t *testing.T) {
	data := []byte(`{
		"focus": [1, 2, 3],
		"top_left": [4, 5, 6],
		"top_right": [7, 8, 9],
		"bottom_left": [10, 11, 12]
	}`)
	c := &PinholeCamera{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		t.Error(err)
	}

	assertEqualVectors(t, maths.NewVec3(1, 2, 3), &c.Focus)
	assertEqualVectors(t, maths.NewVec3(4, 5, 6), &c.TopLeft)
	assertEqualVectors(t, maths.NewVec3(7, 8, 9), &c.TopRight)
	assertEqualVectors(t, maths.NewVec3(10, 11, 12), &c.BottomLeft)
}

func TestAnyCameraJson(t *testing.T) {
	data := []byte(`{
		"type": 		"pinhole",
		"focus":		[0, 0, 0],
		"top_left": 	[-1, 1, 1],
		"top_right":	[1, 1, 1],
		"bottom_left": 	[-1, 1, -1]
	}`)
	c := &AnyCamera{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		t.Error(err)
	}

	ray := c.ShootRay(0.5, 0.5)

	assertEqualVectors(t, maths.NewVec3(0, 1, 0), &ray.Direction)
	assertEqualVectors(t, maths.NewVec3(0, 0, 0), &ray.Start)
}

func assertEqualVectors(t *testing.T, expected *maths.Vec3, v *maths.Vec3) {
	assert := assert.New(t)
	assert.InDelta(expected.X, v.X, maths.Epsilon)
	assert.InDelta(expected.Y, v.Y, maths.Epsilon)
	assert.InDelta(expected.Z, v.Z, maths.Epsilon)
}
