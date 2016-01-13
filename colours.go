package traytor

import "image/color"

//Colour is a representation of a float64 RGB colour
type Colour struct {
	R, G, B float64
}

//Colour32Bit is 32bit colour implementing the color.Color interface
type Colour32Bit struct {
	R, G, B uint32
}

//NewColour32Bit return a new 32bit colour
func NewColour32Bit(r, g, b uint32) *Colour32Bit {
	return &Colour32Bit{R: r, G: g, B: b}
}

//RGBA implements the color.Color interface converting the 32bit colour to 32bit colour with alpha
func (c *Colour32Bit) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, 1
}

//NewColour returns a new RGB colour
func NewColour(r, g, b float64) *Colour {
	return &Colour{R: r, G: g, B: b}
}

//To32Bit returns each of the components of the given RGB color to  uint32
func (c *Colour) To32Bit() *Colour32Bit {
	return NewColour32Bit(uint32(FloatToInt(c.R)), uint32(FloatToInt(c.G)), uint32(FloatToInt(c.B)))
}

//ToColour takes any colour that implements the color.Color interface and turns it into RGB colout(r, g, b are between 0 and 1)
func ToColour(c color.Color) *Colour {
	r, g, b, _ := c.RGBA()
	return NewColour(float64(r)/255.0, float64(g)/255.0, float64(b)/255.0)
}

//MakeZero returns black RGB colour
func (c *Colour) MakeZero() {
	c.SetColour(0, 0, 0)
}

//SetColour sets the colour's components to the given r, g and b
func (c *Colour) SetColour(r, g, b float64) {
	c.R, c.G, c.B = r, g, b
}

//Intensity returns the intensity of the given colour
func (c *Colour) Intensity() float64 {
	return (c.R + c.G + c.B) / 3.0
}

func (c *Colour) Add(other *Colour) {
	c.R += other.R
	c.G += other.G
	c.B += other.B
}
