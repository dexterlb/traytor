package traytor

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"
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

func ExampleImage_Bounds() {
	im := NewImage(640, 480)

	bounds := im.Bounds()
	fmt.Printf("%s\n", bounds)

	// Output:
	// (0,0)-(640,480)
	//
}

func ExampleImage_ColorModel() {
	im := NewImage(640, 480)
	if im.ColorModel() == color.RGBAModel {
		fmt.Printf("Model is RGBA\n")
	}

	// Output:
	// Model is RGBA
	//
}

func ExampleImage_At() {
	im := NewImage(2, 2)
	for j := 0; j < im.height; j++ {
		for i := 0; i < im.width; i++ {
			im.pixels[i][j].SetColour(float32(i), float32(j), 0)
		}
	}
	fmt.Printf("%s\n", im.At(0, 0))
	fmt.Printf("%s\n", im.At(1, 0))
	fmt.Printf("%s\n", im.At(0, 1))
	fmt.Printf("%s\n", im.At(1, 1))

	// Output:
	// [0, 0, 0]
	// [65535, 0, 0]
	// [0, 65535, 0]
	// [65535, 65535, 0]
	//
}

func TestSavePng(t *testing.T) {
	file, err := os.Create("/tmp/test.png")
	if err != nil {
		t.Error(err)
	}

	im := NewImage(20, 20)
	im.pixels[5][5].SetColour(0, 1, 0.5)

	err = png.Encode(file, im)
	if err != nil {
		t.Error(err)
	}
}

/*
	// works, but the test itself is unfinished
func TestLoadPng(t *testing.T) {
	file, err := os.Open("/tmp/test2.png")
	if err != nil {
		t.Error(err)
	}

	decoded, err := png.Decode(file)
	if err != nil {
		t.Error(err)
	}

	im := ToImage(decoded)
	r, g, b, a := decoded.At(5, 5).RGBA()
	t.Error("%s %d %d %d %d\n", im.pixels[5][5], r, g, b, a)
}

*/
