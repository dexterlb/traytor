package materials

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
	"github.com/DexterLB/traytor/sampler"
)

// ReflectiveMaterial is a reflective material
type ReflectiveMaterial struct {
	Colour    *sampler.AnySampler
	Roughness *sampler.AnySampler
}

// Shade returns the emitted colour after intersecting the material
func (m *ReflectiveMaterial) Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour {
	incoming := intersection.Incoming
	colour := m.Colour.GetColour(intersection)
	reflectedRay := &ray.Ray{Depth: incoming.Depth + 1}
	reflectedRay.Direction = *incoming.Direction.Reflected(intersection.Normal)
	reflectedRay.Start = *maths.AddVectors(intersection.Point, intersection.Normal.Scaled(maths.Epsilon))
	return hdrcolour.MultiplyColours(raytracer.Raytrace(reflectedRay), colour)
}
