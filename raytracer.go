package traytor

type Raytracer struct {
	Scene  *Scene
	Random *Random
}

//Raytrace returns the colour obtained by tracing the given ray
func (r *Raytracer) Raytrace(ray *Ray) *Colour {
	intersectionInfo, material := r.Scene.Mesh.Intersect(ray)
	if intersectionInfo == nil {
		return NewColour(0, 0, 0)
	} else {
		return r.Scene.Materials[material].Shade(intersectionInfo, r)
	}
}

//Sample adds another sample to the image by changing it.
func (r *Raytracer) Sample(image *Image) {
	var ray *Ray
	var colour *Colour
	for i := 0; i < image.Width; i++ {
		for j := 0; j < image.Height; j++ {
			ray = r.Scene.Camera.ShootRay(float64(i)/float64(image.Width), float64(j)/float64(image.Height))
			colour = r.Raytrace(ray)
			image.Pixels[i][j].Add(colour)
		}
	}
	image.Divisor += 1
}
