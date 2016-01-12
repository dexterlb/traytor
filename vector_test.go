package traytor

import (
	"fmt"
)

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

func ExampleVec3_Length() {
	v := NewVec3(1, 2, 2)
	fmt.Printf("%.3g\n", v.Length())

	// Output:
	// 3
	//
}
