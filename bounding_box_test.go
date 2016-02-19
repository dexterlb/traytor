package traytor

import "fmt"

func ExampleBoundingBox_AddPoint() {
	box := NewBoundingBox()

	box.AddPoint(NewVec3(-1, -1, -1))
	box.AddPoint(NewVec3(0, 5, -0.5))
	box.AddPoint(NewVec3(1, 0, 2))

	fmt.Printf("min: %s, max: %s\n", box.MinVolume, box.MaxVolume)

	// Output:
	// min: (-1, -1, -1), max: (1, 5, 2)
}

func ExampleBoundingBox_Inside() {
	box := &BoundingBox{
		MinVolume: NewVec3(-1, -1, -1),
		MaxVolume: NewVec3(1, 1, 1),
	}

	fmt.Printf("0, 0, 0: %v\n", box.Inside(NewVec3(0, 0, 0)))
	fmt.Printf("0, 0, 2: %v\n", box.Inside(NewVec3(0, 0, 2)))

	// Output:
	// 0, 0, 0: true
	// 0, 0, 2: false
}

func ExampleBoundingBox_Intersect() {
	box := &BoundingBox{
		MinVolume: NewVec3(-1, -1, -1),
		MaxVolume: NewVec3(1, 1, 1),
	}

	ray1 := &Ray{
		Start:     *NewVec3(0, 0, 0),
		Direction: *NewVec3(1, 0, 0),
	}

	ray2 := &Ray{
		Start:     *NewVec3(-5, 0, 0.5),
		Direction: *NewVec3(1, 0, 0),
	}

	ray3 := &Ray{
		Start:     *NewVec3(-5, 0, 0.5),
		Direction: *NewVec3(-1, 0, 0),
	}

	ray4 := &Ray{
		Start:     *NewVec3(-5, 0, 0),
		Direction: *NewVec3(1, 0, 5),
	}

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
