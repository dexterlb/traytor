package traytor

import (
	"image"
	"image/color"
)

//Image is a stuct which will display images via its 2D colour array, wich represents the screen
type Image struct {
	Pixels        [][]Colour
	Width, Height int
	Divisor       int
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
	return &Image{Pixels: pixels, Width: width, Height: height, Divisor: 1}
}

//String returns a string which is the representaton of image: {r, g, b}, ... {r, g, b}\n ...\n {r, g, b},...{r, g, b}
func (im *Image) String() string {
	representation := ""
	for j := 0; j < im.Width; j++ {
		for i := 0; i < im.Height; i++ {
			representation += im.Pixels[i][j].String()
			if i != im.Width-1 {
				representation += ", "
			}
		}
		representation += "\n"
	}
	return representation
}

//Add adds another image to this one
func (im *Image) Add(other *Image) {
	for j := 0; j < im.Width; j++ {
		for i := 0; i < im.Height; i++ {
			(&im.Pixels[i][j]).Add(&other.Pixels[i][j])
		}
	}
}

//AddImages returns a new image which is the sum of the given ones
func AddImages(a *Image, b *Image) *Image {
	sum := NewImage(a.Width, a.Height)
	sum.Add(a)
	sum.Add(b)
	return sum
}

//AtHDR returns the Colour of the pixel at [x][y] (scaled by the divisor)
func (im *Image) AtHDR(x, y int) *Colour {
	if im.Divisor == 0 {
		return NewColour(1, 1, 1)
	}
	return im.Pixels[x][y].Scaled(1 / float32(im.Divisor))
}

//At returns the sRGB Colour of the pixel at [x][y (scaled by the divisor)]
func (im *Image) At(x, y int) color.Color {
	return im.AtHDR(x, y).To32Bit()
}

//ColorModel returns the image's color model (as used by Go's image interface)
func (im *Image) ColorModel() color.Model {
	return color.RGBAModel
}

//Bounds returns a rectangle as big as the image
func (im *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.Width, im.Height)
}

//ToImage constructs an Image from an sRGB image
func ToImage(im image.Image) *Image {
	width := im.Bounds().Max.X - im.Bounds().Min.X
	height := im.Bounds().Max.Y - im.Bounds().Min.Y
	extractedImage := NewImage(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			extractedImage.Pixels[i][j].Add(ToColour(im.At(im.Bounds().Min.X+i, im.Bounds().Min.Y+j)))
		}
	}
	return extractedImage
}
