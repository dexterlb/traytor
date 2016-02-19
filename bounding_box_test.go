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
