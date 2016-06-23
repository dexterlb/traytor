package materials

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/ray"
	"github.com/DexterLB/traytor/sampler"
)

// EmissiveMaterial acts as a lamp
type EmissiveMaterial struct {
	Colour   *sampler.AnySampler
	Strength *sampler.AnySampler
}

// Shade returns the emitted colour after intersecting the material
func (m *EmissiveMaterial) Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour {
	return m.Colour.GetColour(intersection).Scaled(float32(m.Strength.GetFac(intersection)))
}
