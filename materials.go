package traytor

import (
	"encoding/json"
	"fmt"
)

// AnyMaterial implements the Material interface and is deserialiseable from json
type AnyMaterial struct {
	Material
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *AnyMaterial) UnmarshalJson(data []byte) error {
	materialType, err := jsonObjectType(data)
	if err != nil {
		return err
	}

	// Here comes code duplication
	switch materialType {
	case "emissive":
		material := &EmissiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "lambert":
		material := &LambertMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "reflective":
		material := &ReflectiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "refractive":
		material := &RefractiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	default:
		return fmt.Errorf("Unknown material type: '%s'", materialType)
	}

	return nil
}

// Material objects are used to shade surfaces
type Material interface {
	Shade(intersection *Intersection, raytracer *Raytracer) *Colour
}

// EmissiveMaterial acts as a lamp
type EmissiveMaterial struct {
	Colour   *Colour
	Strength float32
}

// Shade returns the emitted colour after intersecting the material
func (m *EmissiveMaterial) Shade(intersection *Intersection, raytracer *Raytracer) *Colour {
	return m.Colour.Scaled(m.Strength)
}

// LambertMaterial is a simple diffuse material
type LambertMaterial struct {
	Colour *Colour
}

// Shade returns the emitted colour after intersecting the material
func (m *LambertMaterial) Shade(intersection *Intersection, raytracer *Raytracer) *Colour {
	return NewColour(0, 0, 0)
}

// ReflectiveMaterial is a reflective material
type ReflectiveMaterial struct {
	Colour    *Colour
	Roughness float32
}

// Shade returns the emitted colour after intersecting the material
func (m *ReflectiveMaterial) Shade(intersection *Intersection, raytracer *Raytracer) *Colour {
	return NewColour(0, 0, 0)
}

// RefractiveMaterial is a material for modeling glass, etc
type RefractiveMaterial struct {
	Colour    *Colour
	Roughness float32
	IOR       float32
}

// Shade returns the emitted colour after intersecting the material
func (m *RefractiveMaterial) Shade(intersection *Intersection, raytracer *Raytracer) *Colour {
	return NewColour(0, 0, 0)
}
