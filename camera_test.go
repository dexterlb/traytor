package traytor

import "fmt"

func ExamplePinholeCamera_ShootRay() {
	var c Camera = &PinholeCamera{
		Focus:      *NewVec3(0, 0, 0),
		TopLeft:    *NewVec3(-1, 1, 1),
		TopRight:   *NewVec3(1, 1, 1),
		BottomLeft: *NewVec3(-1, 1, -1),
	}
	ray := c.ShootRay(0, 0)
	fmt.Printf("%s -> %s\n", &ray.Start, &ray.Direction)
	// Output:
	// (0, 0, 0) -> (-1, 1, 1)
	//
}
