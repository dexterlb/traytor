package materials

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
	"github.com/DexterLB/traytor/sampler"
)

// LambertMaterial is a simple diffuse material
type LambertMaterial struct {
	Colour *sampler.AnySampler
}

// Shade returns the emitted colour after intersecting the material
func (m *LambertMaterial) Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour {
	randomRayDir := *raytracer.RandomGen().Vec3HemiCos(intersection.Normal)
	randomRayStart := *maths.AddVectors(intersection.Point, intersection.Normal.Scaled(maths.Epsilon))
	ray := &ray.Ray{Start: randomRayStart, Direction: randomRayDir, Depth: intersection.Incoming.Depth + 1}
	colour := raytracer.Raytrace(ray)
	colour.MultiplyBy(m.Colour.GetColour(intersection))
	return colour
}
