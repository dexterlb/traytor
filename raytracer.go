package traytor

type Raytracer struct {
	scene  *Scene
	random *Random
}

//Raytrace returns the colour obtained by tracing the given ray
func (r *Raytracer) Raytrace(ray *Ray) *Colour {
	return NewColour(0, 0, 0)
}

//Sample adds another sample to the image by changing it.
func (r *Raytracer) Sample(image *Image) {

}

//LoadScene load the scene from a file
func LoadScene(filename string) *Image {
	return NewImage(400, 400)
}
