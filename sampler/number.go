package sampler

import (
	"encoding/json"

	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// NumberSampler is a sampler constructed by a single number
type NumberSampler struct {
	Colour *hdrcolour.Colour
	Vector *maths.Vec3
	Fac    float64
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (n *NumberSampler) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &n.Fac)
	if err != nil {
		return err
	}
	n.Colour = hdrcolour.New(
		float32(n.Fac),
		float32(n.Fac),
		float32(n.Fac),
	)
	n.Vector = maths.NewVec3(n.Fac, n.Fac, n.Fac)

	return nil
}

// GetColour returns a colour with R, G, B equal to the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetColour(intersection *ray.Intersection) *hdrcolour.Colour {
	return n.Colour
}

// GetVec3 returns a vector with X, Y, Z equal to the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetVec3(intersection *ray.Intersection) *maths.Vec3 {
	return n.Vector
}

// GetFac returns the number
// (doesn't need a valid intersection)
func (n *NumberSampler) GetFac(intersection *ray.Intersection) float64 {
	return n.Fac
}
