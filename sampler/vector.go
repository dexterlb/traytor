package sampler

import (
	"encoding/json"

	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// Vec3Sampler is a sampler constructed by 3 numbers
type Vec3Sampler struct {
	Colour *hdrcolour.Colour
	Vector *maths.Vec3
	Fac    float64
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Vec3Sampler) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &v.Vector)
	if err != nil {
		return err
	}
	v.Colour = hdrcolour.New(
		float32(v.Vector.X),
		float32(v.Vector.Y),
		float32(v.Vector.Z),
	)
	v.Fac = (v.Vector.X + v.Vector.Y + v.Vector.Z) / 3 // fixme: this must be smarter

	return nil
}

// GetColour returns a colour with R, G, B corresponding to the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetColour(intersection *ray.Intersection) *hdrcolour.Colour {
	return v.Colour
}

// GetVec3 returns a vector with X, Y, Z corresponding to the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetVec3(intersection *ray.Intersection) *maths.Vec3 {
	return v.Vector
}

// GetFac returns a number that is the average of the 3 components
// (doesn't need a valid intersection)
func (v *Vec3Sampler) GetFac(intersection *ray.Intersection) float64 {
	return v.Fac
}
