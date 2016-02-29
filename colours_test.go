package traytor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualColours(t *testing.T, expected *Colour, c *Colour) {
	assert := assert.New(t)
	colourEpsilon := 0.0001

	assert.InDelta(expected.R, c.R, colourEpsilon)
	assert.InDelta(expected.G, c.G, colourEpsilon)
	assert.InDelta(expected.B, c.B, colourEpsilon)
}

func TestTo32Bit(t *testing.T) {
	assert := assert.New(t)
	c := NewColour(1, 1, 1)
	intColour := c.To32Bit()
	assert.Equal(uint32(65535), intColour.R)
	assert.Equal(uint32(65535), intColour.G)
	assert.Equal(uint32(65535), intColour.B)

	c = NewColour(0, 0, 0)
	intColour = c.To32Bit()
	assert.Equal(uint32(0), intColour.R)
	assert.Equal(uint32(0), intColour.G)
	assert.Equal(uint32(0), intColour.B)

	c = NewColour(0.000001, 0, 0)
	intColour = c.To32Bit()
	assert.Equal(uint32(1), intColour.R)
	assert.Equal(uint32(0), intColour.G)
	assert.Equal(uint32(0), intColour.B)

	c = NewColour(0.7, 0.2, 0.5)
	intColour = c.To32Bit()
	assert.Equal(uint32(0xdab3), intColour.R)
	assert.Equal(uint32(0x7c0a), intColour.G)
	assert.Equal(uint32(0xbc40), intColour.B)
}

func TestRGBA(t *testing.T) {
	assert := assert.New(t)
	c := NewColour32Bit(51, 0, 0)
	r, g, b, a := c.RGBA()
	assert.Equal(uint32(51), r)
	assert.Equal(uint32(0), g)
	assert.Equal(uint32(0), b)
	assert.Equal(uint32(65535), a)
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

func TestToColour(t *testing.T) {
	c := NewColour32Bit(0, 0, 0)
	assertEqualColours(t, NewColour(0, 0, 0), ToColour(c))

	c = NewColour32Bit(0xdab3, 0x7c0a, 0xbc40)
	assertEqualColours(t, NewColour(0.7, 0.2, 0.5), ToColour(c))
}

func TestColourJson(t *testing.T) {
	c := NewColour(0, 0, 0)
	err := json.Unmarshal([]byte(`[0.4, 0.5, 1]`), &c)
	if err != nil {
		t.Error(err)
	}
	assertEqualColours(t, NewColour(0.4, 0.5, 1), c)
}
