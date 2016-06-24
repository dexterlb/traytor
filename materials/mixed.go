package materials

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/ray"
)

// MixedMaterials mixes the shaders of two materials where
type MixedMaterial struct {
	First       *AnyMaterial
	Second      *AnyMaterial
	Coefficient float64
}

// Shade returns colour depending on coefficient and random number
func (m *MixedMaterial) Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour {
	if raytracer.RandomGen().Float01() < m.Coefficient {
		return m.First.Shade(intersection, raytracer)
	}
	return m.Second.Shade(intersection, raytracer)
}
