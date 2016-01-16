package traytor

import (
	"image"
	"image/color"
)

//Image is a stuct which will display images via its 2D colour array, wich represents the screen
type Image struct {
	pixels        [][]Colour
	width, height int
}

//NewImage will set the screen to the given width and height
func NewImage(width, height int) *Image {
	pixels := make([][]Colour, width)
	for i := range pixels {
		pixels[i] = make([]Colour, height)
		for j := range pixels[i] {
			pixels[i][j] = *NewColour(0, 0, 0)
		}
	}
	return &Image{pixels: pixels, width: width, height: height}
}

//String returns a string which is the representaton of image: {r, g, b}, ... {r, g, b}\n ...\n {r, g, b},...{r, g, b}
func (im *Image) String() string {
	representation := ""
	for j := 0; j < im.height; j++ {
		for i := 0; i < im.width; i++ {
			representation += im.pixels[i][j].String()
			if i != im.width-1 {
				representation += ", "
			}
		}
		representation += "\n"
	}
	return representation
}

//Add returns a new image which is the sum of two given.
func (im *Image) Add(other *Image) *Image {
	sum := NewImage(im.width, im.height)
	for j := 0; j < im.height; j++ {
		for i := 0; i < im.width; i++ {
			sum.pixels[i][j] = *Sum(im.pixels[i][j], other.pixels[i][j])
		}
	}
	return sum
}

//At returns the Colour of the pixel at [x][y]
func (im *Image) At(x, y int) color.Color {
	return im.pixels[x][y].To32Bit()
}

//ColorModel returns the image's color model (as used by Go's image interface)
func (im *Image) ColorModel() color.Model {
	return color.RGBAModel
}

//Bounds returns a rectangle as big as the image
func (im *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.width, im.height)
}
