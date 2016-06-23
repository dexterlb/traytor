package raytracer

import (
	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/random"
	"github.com/DexterLB/traytor/ray"
	"github.com/DexterLB/traytor/scene"
)

// Raytracer represents a single rendering unit
type Raytracer struct {
	Scene  *scene.Scene
	Random *random.Random
}

// RandomGen returns the raytracer's random generator
func (r *Raytracer) RandomGen() *random.Random {
	return r.Random
}

// Raytrace returns the colour obtained by tracing the given ray
func (r *Raytracer) Raytrace(incoming *ray.Ray) *hdrcolour.Colour {
	if incoming.Depth > r.Scene.MaxDepth {
		return hdrcolour.New(0, 0, 0)
	}
	intersectionInfo := r.Scene.Mesh.Intersect(incoming)
	if intersectionInfo == nil {
		return hdrcolour.New(0, 0, 0)
	}
	return r.Scene.Materials[intersectionInfo.Material].Shade(intersectionInfo, r)
}

// Sample adds another sample to the image by changing it.
func (r *Raytracer) Sample(image *hdrimage.Image) {
	var ray *ray.Ray
	var colour *hdrcolour.Colour
	for i := 0; i < image.Width; i++ {
		for j := 0; j < image.Height; j++ {
			ray = r.Scene.Camera.ShootRay(
				(float64(i)+r.Random.Float01())/float64(image.Width),
				(float64(j)+r.Random.Float01())/float64(image.Height),
			)
			colour = r.Raytrace(ray)
			image.Pixels[i][j].Add(colour)
		}
	}
	image.Divisor++
}
