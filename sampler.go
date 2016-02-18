package traytor

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

// AnySampler implements the Sampler interface and is deserialiseable from json
type AnySampler struct {
	Sampler
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *AnySampler) UnmarshalJSON(data []byte) error {
	numberSampler := &NumberSampler{}
	err := json.Unmarshal(data, &numberSampler)
	if err == nil {
		*s = AnySampler{numberSampler}
		return nil
	}

	vec3Sampler := &Vec3Sampler{}
	err = json.Unmarshal(data, &vec3Sampler)
	if err == nil {
		*s = AnySampler{vec3Sampler}
		return nil
	}

	samplerType, err := jsonObjectType(data)
	if err != nil {
		return err
	}

	switch samplerType {
	case "image_texture":
		sampler := &ImageTexture{}
		err = json.Unmarshal(data, &sampler)
		if err != nil {
			return err
		}
		*s = AnySampler{sampler}
	default:
		return fmt.Errorf("Unknown texture sampler: '%s'", samplerType)
	}

	return nil
}

// Sampler objects return varying attributes on points on surfaces
// (e.g. texture colours, ray information etc)
type Sampler interface {
	GetColour(intersection *Intersection) *Colour
	GetVec3(intersection *Intersection) *Vec3
	GetFac(intersectuib *Intersection) float64
}

// Vec3Sampler is a sampler constructed by 3 numbers
type Vec3Sampler struct {
	Colour *Colour
	Vector *Vec3
	Fac    float64
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Vec3Sampler) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &v.Vector)
	if err != nil {
		return err
	}
	v.Colour = NewColour(
		float32(v.Vector.X),
		float32(v.Vector.Y),
		float32(v.Vector.Z),
	)
	v.Fac = (v.Vector.X + v.Vector.Y + v.Vector.Z) / 3 // fixme: this must be smarter

	return nil
}

// GetColour returns a colour with R, G, B corresponding to the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetColour(intersection *Intersection) *Colour {
	return v.Colour
}

// GetVec3 returns a vector with X, Y, Z corresponding to the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetVec3(intersection *Intersection) *Vec3 {
	return v.Vector
}

// GetFac returns a number that is the average of the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetFac(intersection *Intersection) float64 {
	return v.Fac
}

// NumberSampler is a sampler constructed by a single number
type NumberSampler struct {
	Colour *Colour
	Vector *Vec3
	Fac    float64
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (n *NumberSampler) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &n.Fac)
	if err != nil {
		return err
	}
	n.Colour = NewColour(
		float32(n.Fac),
		float32(n.Fac),
		float32(n.Fac),
	)
	n.Vector = NewVec3(n.Fac, n.Fac, n.Fac)

	return nil
}

// GetColour returns a colour with R, G, B equal to the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetColour(intersection *Intersection) *Colour {
	return n.Colour
}

// GetVec3 returns a vector with X, Y, Z equal to the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetVec3(intersection *Intersection) *Vec3 {
	return n.Vector
}

// GetFac returns the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetFac(intersection *Intersection) float64 {
	return n.Fac
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *ImageTexture) UnmarshalJSON(data []byte) error {
	textureSettings := &struct {
		Scale    *Vec3
		Offset   *Vec3
		Filename string
	}{}

	err := json.Unmarshal(data, textureSettings)
	if err != nil {
		return err
	}

	f, err := os.Open(textureSettings.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fileType := filepath.Ext(textureSettings.Filename)

	var ldrImage image.Image

	switch fileType {
	case ".png":
		ldrImage, err = png.Decode(f)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown filetype for image texture: %s\n", fileType)
	}

	i.Image = ToImage(ldrImage)

	i.ScaleU = textureSettings.Scale.X
	i.ScaleV = textureSettings.Scale.Y

	i.OffsetU = textureSettings.Offset.X
	i.OffsetV = textureSettings.Offset.Y

	return nil
}

// ImageTexture is.. an image texture!
type ImageTexture struct {
	Image                            *Image
	ScaleU, ScaleV, OffsetU, OffsetV float64
}

// GetColour returns the colour at (u, v)
func (i *ImageTexture) GetColour(intersection *Intersection) *Colour {
	u := i.OffsetU + i.ScaleU*intersection.U
	v := i.OffsetV + i.ScaleV*intersection.V

	// tile it
	u -= math.Floor(u)
	v -= math.Floor(v)

	return i.Image.AtHDR(int(u*float64(i.Image.Width)), int(v*float64(i.Image.Height)))
}

// GetFac returns the colour intensity at (u, v)
func (i *ImageTexture) GetFac(intersection *Intersection) float64 {
	return float64(i.GetColour(intersection).Intensity())
}

// GetVec3 returns the colour at (u, v)'s components as a vector
func (i *ImageTexture) GetVec3(intersection *Intersection) *Vec3 {
	colour := i.GetColour(intersection)
	return NewVec3(float64(colour.R), float64(colour.G), float64(colour.B))
}
