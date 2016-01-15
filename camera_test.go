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
	fmt.Printf("0, 0: %s -> %s\n", &ray.Start, &ray.Direction)

	ray = c.ShootRay(0.5, 0.5)
	fmt.Printf("0.5, 0.5: %s -> %s\n", &ray.Start, &ray.Direction)

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
	fmt.Printf("0, 0: %s -> %s\n", &ray.Start, &ray.Direction)

	ray = c.ShootRay(0.5, 0.5)
	fmt.Printf("0.5, 0.5: %s -> %s\n", &ray.Start, &ray.Direction)

	// Output:
	// 0, 0: (-2, 15, 3) -> (-0.577, 0.577, 0.577)
	// 0.5, 0.5: (-2, 15, 3) -> (0, 1, 0)
}
