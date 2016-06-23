package materials

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
	"github.com/DexterLB/traytor/sampler"
)

// RefractiveMaterial is a material for modeling glass, etc
type RefractiveMaterial struct {
	Colour    *sampler.AnySampler
	Roughness *sampler.AnySampler
	IOR       *sampler.AnySampler
}

// Shade returns the emitted colour after intersecting the material
func (m *RefractiveMaterial) Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour {
	incoming := &intersection.Incoming.Direction
	normal := intersection.Normal
	ior := m.IOR.GetFac(intersection)
	colour := m.Colour.GetColour(intersection)
	refracted := &maths.Vec3{}
	startPoint := &maths.Vec3{}

	if maths.DotProduct(incoming, normal) < 0 {
		refracted = maths.Refract(incoming, normal, 1/ior)
	} else {
		refracted = maths.Refract(incoming, normal.Negative(), ior)
	}

	if refracted != nil {
		// regular refraction - push the starting point a tiny bit
		// through the surface
		startPoint = maths.MinusVectors(
			intersection.Point, normal.FaceForward(incoming).Scaled(maths.Epsilon),
		)
	} else {
		// total inner reflection
		refracted = incoming.Reflected(normal.FaceForward(incoming))
		// push the starting point a tiny bit away from the surface
		startPoint = maths.AddVectors(
			intersection.Point, normal.FaceForward(incoming).Scaled(maths.Epsilon),
		)
	}

	newRay := &ray.Ray{Depth: intersection.Incoming.Depth + 1}
	newRay.Start = *startPoint
	newRay.Direction = *refracted
	return hdrcolour.MultiplyColours(raytracer.Raytrace(newRay), colour)

}
