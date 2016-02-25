package traytor

import (
	"encoding/json"
	"fmt"
	"image/color"
)

//Colour is a representation of a float32 RGB colour
type Colour struct {
	R, G, B float32
}

//String returns the string representation of the colour in the form of {r, g, b}
func (c *Colour) String() string {
	return fmt.Sprintf("{%.3g, %.3g, %.3g}", c.R, c.G, c.B)
}

//String returns the string representation of the 32bit colour in the form of [r, g, b]
func (c *Colour32Bit) String() string {
	return fmt.Sprintf("[%d, %d, %d]", c.R, c.G, c.B)
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
	return c.R, c.G, c.B, 65535
}

//NewColour returns a new RGB colour
func NewColour(r, g, b float32) *Colour {
	return &Colour{R: r, G: g, B: b}
}

//To32Bit returns each of the components of the given RGB color to  uint32
func (c *Colour) To32Bit() *Colour32Bit {
	return NewColour32Bit(linearTosRGB(c.R), linearTosRGB(c.G), linearTosRGB(c.B))
}

//linearTosRGBreturn an int between 0 and 1 constructed from a given float between 0 and 65535
func linearTosRGB(x float32) uint32 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 65535
	}
	if x <= 0.00313008 {
		x = x * 12.02
	} else {
		x = (1.055)*Pow32(x, 1.0/2.4) - 0.055
	}
	return uint32(Round32(x * 65535.0))
}

//sRGBToLinear converts singel int number to float using special magic formula.
func sRGBToLinear(i uint32) float32 {
	if i > 65535 {
		return 1
	}

	x := float32(i) / 65535.0
	if x <= 0.04045 {
		return x / 12.92
	}
	return Pow32((x+0.055)/1.055, 2.4)
}

//ToColour takes any colour that implements the color.Color interface and turns it into RGB colout(r, g, b are between 0 and 1)
func ToColour(c color.Color) *Colour {
	r, g, b, _ := c.RGBA()
	return NewColour(sRGBToLinear(r), sRGBToLinear(g), sRGBToLinear(b))
}

//MakeZero returns black RGB colour
func (c *Colour) MakeZero() {
	c.SetColour(0, 0, 0)
}

//SetColour sets the colour's components to the given r, g and b
func (c *Colour) SetColour(r, g, b float32) {
	c.R, c.G, c.B = r, g, b
}

//Intensity returns the intensity of the given colour
func (c *Colour) Intensity() float32 {
	return (c.R + c.G + c.B) / 3.0
}

//Add adds another colour to this one
func (c *Colour) Add(other *Colour) {
	c.R += other.R
	c.G += other.G
	c.B += other.B
}

//Scale multiplies the colour by the given multiplier
func (c *Colour) Scale(multiplier float32) {
	c.R *= multiplier
	c.G *= multiplier
	c.B *= multiplier
}

//Scaled returns a new colour which is the product of the original and multiplier
func (c *Colour) Scaled(multiplier float32) *Colour {
	return NewColour(
		c.R*multiplier,
		c.G*multiplier,
		c.B*multiplier,
	)
}

//UnmarshalJSON implements the json.Unmarshaler interface
func (c *Colour) UnmarshalJSON(data []byte) error {
	var unmarshaled []float32
	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}
	c.R = unmarshaled[0]
	c.G = unmarshaled[1]
	c.B = unmarshaled[2]
	return nil
}

//AddColours adds two colours
func AddColours(first, other *Colour) *Colour {
	r := first.R + other.R
	g := first.G + other.G
	b := first.B + other.B
	return NewColour(r, g, b)
}

//MultiplyColours returns the Product of two colours
func MultiplyColours(first, other *Colour) *Colour {
	r := first.R * other.R
	g := first.G * other.G
	b := first.B * other.B
	return NewColour(r, g, b)
}

//MultiplyBy other colour sets the given vector to its product with the other colour
func (c *Colour) MultiplyBy(other *Colour) {
	c.R *= other.R
	c.G *= other.G
	c.B *= other.B
}
