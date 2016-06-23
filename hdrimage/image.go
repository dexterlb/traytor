package hdrimage

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/DexterLB/traytor/hdrcolour"
)

// Image is a stuct which will display images via its 2D colour array, wich represents the screen
type Image struct {
	Pixels        [][]hdrcolour.Colour
	Width, Height int
	Divisor       int
}

// New will set the screen to the given width and height
func New(width, height int) *Image {
	pixels := make([][]hdrcolour.Colour, width)
	for i := range pixels {
		pixels[i] = make([]hdrcolour.Colour, height)
		for j := range pixels[i] {
			pixels[i][j] = *hdrcolour.New(0, 0, 0)
		}
	}
	return &Image{Pixels: pixels, Width: width, Height: height, Divisor: 1}
}

// Decode reads data in the simple traytor_hdr format and produces an
// image.
func Decode(reader io.Reader) (*Image, error) {
	size := [2]uint16{}
	err := binary.Read(reader, binary.LittleEndian, &size)
	if err != nil {
		return nil, fmt.Errorf("cannot read header from image data: %s", err)
	}

	im := New(int(size[0]), int(size[1]))
	var rgba [4]float32

	for i := im.Width - 1; i > 0; i-- {
		for j := 0; j < im.Height; j++ {
			err = binary.Read(reader, binary.LittleEndian, &rgba)
			if err != nil {
				return nil, fmt.Errorf("cannot read image data: %s", err)
			}
			im.Pixels[j][i].R = rgba[0]
			im.Pixels[j][i].G = rgba[1]
			im.Pixels[j][i].B = rgba[2]
		}
	}

	return im, nil
}

// String returns a string which is the representaton of image:
// {r, g, b}, ... {r, g, b}\n ...\n {r, g, b},...{r, g, b}
func (im *Image) String() string {
	representation := ""
	for i := 0; i < im.Width; i++ {
		for j := 0; j < im.Height; j++ {
			representation += im.Pixels[i][j].String()
			if j != im.Width-1 {
				representation += ", "
			}
		}
		representation += "\n"
	}
	return representation
}

// Add adds another image to this one
func (im *Image) Add(other *Image) {
	for i := 0; i < im.Width; i++ {
		for j := 0; j < im.Height; j++ {
			(&im.Pixels[i][j]).Add(&other.Pixels[i][j])
		}
	}
}

// Add returns a new image which is the sum of the given ones
func Add(a *Image, b *Image) *Image {
	sum := New(a.Width, a.Height)
	sum.Add(a)
	sum.Add(b)
	return sum
}

// AtHDR returns the Colour of the pixel at [x][y] (scaled by the divisor)
func (im *Image) AtHDR(x, y int) *hdrcolour.Colour {
	if im.Divisor == 0 {
		return hdrcolour.New(1, 1, 1)
	}
	return im.Pixels[x][y].Scaled(1 / float32(im.Divisor))
}

// At returns the sRGB Colour of the pixel at [x][y (scaled by the divisor)]
func (im *Image) At(x, y int) color.Color {
	return im.AtHDR(x, y).To32Bit()
}

// ColorModel returns the image's color model (as used by Go's image interface)
func (im *Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns a rectangle as big as the image
func (im *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.Width, im.Height)
}

// FromSRGB constructs an Image from an sRGB image
func FromSRGB(im image.Image) *Image {
	width := im.Bounds().Max.X - im.Bounds().Min.X
	height := im.Bounds().Max.Y - im.Bounds().Min.Y
	extractedImage := New(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			extractedImage.Pixels[i][j].Add(
				hdrcolour.FromColor(im.At(im.Bounds().Min.X+i, im.Bounds().Min.Y+j)),
			)
		}
	}
	return extractedImage
}
