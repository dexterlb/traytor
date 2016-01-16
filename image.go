package traytor

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

func (im *Image) Add(other *Image) *Image {
	sum := NewImage(im.width, im.height)
	for j := 0; j < im.height; j++ {
		for i := 0; i < im.width; i++ {
			sum.pixels[i][j] = *Sum(im.pixels[i][j], other.pixels[i][j])
		}
	}
	return sum
}
