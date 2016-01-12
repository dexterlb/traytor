package traytor

import "fmt"

func ExampleVec3_ToZero() {
	v := NewVec3(1, 2, 3)
	fmt.Printf("%s\n", v)
	v.ToZero()
	fmt.Printf("%s\n", v)

	// Output:
	// (1, 2, 3)
	// (0, 0, 0)
	//
}
