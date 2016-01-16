package traytor

import (
	"fmt"
)

func ExampleImage_String() {
	i := NewImage(2, 2)
	fmt.Printf("%s\n", i)
	// Output:
	// {0, 0, 0}, {0, 0, 0}
	// {0, 0, 0}, {0, 0, 0}
	//
}

func ExampleImage_Add() {
	im1 := NewImage(2, 2)
	for j := 0; j < im1.height; j++ {
		for i := 0; i < im1.width; i++ {
			im1.pixels[i][j].SetColour(4, 5, 6)
		}
	}
	im2 := NewImage(2, 2)
	for j := 0; j < im2.height; j++ {
		for i := 0; i < im2.width; i++ {
			im2.pixels[i][j].SetColour(1, 2, 3)
		}
	}
	im3 := im1.Add(im2)
	fmt.Printf("%s\n", im3)

	// Output:
	// {5, 7, 9}, {5, 7, 9}
	// {5, 7, 9}, {5, 7, 9}
	//

}
