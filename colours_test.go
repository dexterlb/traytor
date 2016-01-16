package traytor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func asserEqualColours(t *testing.T, expected *Colour, c *Colour) {
	assert := assert.New(t)
	assert.InEpsilon(expected.R, c.R, Epsilon)
	assert.InEpsilon(expected.G, c.G, Epsilon)
	assert.InEpsilon(expected.B, c.B, Epsilon)
}

func TestTo32Bit(t *testing.T) {
	assert := assert.New(t)
	c := NewColour(1, 1, 1)
	intColour := c.To32Bit()
	assert.InEpsilon(255, intColour.R, Epsilon)
	assert.InEpsilon(255, intColour.G, Epsilon)
	assert.InEpsilon(255, intColour.B, Epsilon)

	c = NewColour(0, 0, 0)
	intColour = c.To32Bit()
	assert.InEpsilon(0, intColour.R, Epsilon)
	assert.InEpsilon(0, intColour.G, Epsilon)
	assert.InEpsilon(0, intColour.B, Epsilon)

	c = NewColour(0.0002, 0, 0)
	intColour = c.To32Bit()
	assert.InEpsilon(1, intColour.R, Epsilon)
	assert.InEpsilon(0, intColour.G, Epsilon)
	assert.InEpsilon(0, intColour.B, Epsilon)
}

func TestRGBA(t *testing.T) {
	assert := assert.New(t)
	c := NewColour32Bit(51, 0, 0)
	r, g, b, a := c.RGBA()
	assert.Equal(uint32(51), r)
	assert.Equal(uint32(0), g)
	assert.Equal(uint32(0), b)
	assert.Equal(uint32(255), a)
}

func ExampleColour_String() {
	c := NewColour(0, 0, 0)
	fmt.Printf("%s\n", c)
	c = NewColour(42, 3.04, -12.4)
	fmt.Printf("%s\n", c)
	// Output:
	// {0, 0, 0}
	// {42, 3.04, -12.4}
	//
}

/*
func TestToColour(t *testing.T) {
	c := NewColour32Bit(0, 0, 0)
	asserEqualColours(t, NewColour(0, 0, 0), ToColour(c))

	c = NewColour32Bit(51, 0, 0)
	asserEqualColours(t, NewColour(0.4679389891357439, 0, 0), ToColour(c))
}
*/
