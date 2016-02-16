package traytor

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ExamplePinholeCamera_ShootRay() {
	var c Camera = &PinholeCamera{
		Focus:      *NewVec3(0, 0, 0),
		TopLeft:    *NewVec3(-1, 1, 1),
		TopRight:   *NewVec3(1, 1, 1),
		BottomLeft: *NewVec3(-1, 1, -1),
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
		Focus:      *NewVec3(-2, 15, 3),
		TopLeft:    *NewVec3(-3, 16, 4),
		TopRight:   *NewVec3(-1, 16, 4),
		BottomLeft: *NewVec3(-3, 16, 2),
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

	asserEqualVectors(t, NewVec3(1, 2, 3), &c.Focus)
	asserEqualVectors(t, NewVec3(4, 5, 6), &c.TopLeft)
	asserEqualVectors(t, NewVec3(7, 8, 9), &c.TopRight)
	asserEqualVectors(t, NewVec3(10, 11, 12), &c.BottomLeft)
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

	asserEqualVectors(t, NewVec3(0, 1, 0), &ray.Direction)
	asserEqualVectors(t, NewVec3(0, 0, 0), &ray.Start)
}
