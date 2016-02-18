package traytor

import "encoding/json"

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

	return err
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
