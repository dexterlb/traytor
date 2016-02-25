package traytor

// Raytracer represents a single rendering unit
type Raytracer struct {
	Scene    *Scene
	Random   *Random
	MaxDepth int
}

// Raytrace returns the colour obtained by tracing the given ray
func (r *Raytracer) Raytrace(ray *Ray) *Colour {
	if ray.Depth > r.MaxDepth {
		return NewColour(0, 0, 0)
	}
	intersectionInfo := r.Scene.Mesh.Intersect(ray)
	if intersectionInfo == nil {
		return NewColour(0, 0, 0)
	}
	return r.Scene.Materials[intersectionInfo.Material].Shade(intersectionInfo, r)
}

// Sample adds another sample to the image by changing it.
func (r *Raytracer) Sample(image *Image) {
	var ray *Ray
	var colour *Colour
	for i := 0; i < image.Width; i++ {
		for j := 0; j < image.Height; j++ {
			ray = r.Scene.Camera.ShootRay((float64(i)+r.Random.Float01())/float64(image.Width), (float64(j)+r.Random.Float01())/float64(image.Height))
			colour = r.Raytrace(ray)
			image.Pixels[i][j].Add(colour)
		}
	}
	image.Divisor++
}
