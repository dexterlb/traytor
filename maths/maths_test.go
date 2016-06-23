package maths

import "fmt"

func ExampleSolveEquation() {
	a := NewVec3(2, 1, 0)
	b := NewVec3(4, -2, 0)
	c := NewVec3(24, 4, 0)

	// | 2x + 4y = 24
	// |  x - 2y = 4

	x, y := SolveEquation(a, b, c)
	fmt.Printf("x = %.3g, y = %.3g\n", x, y)

	// Output:
	// x = 8, y = 2
}
