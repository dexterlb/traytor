package traytor

// Camera is a pinhole camera which can shoot rays. It has a focus
// (the location of the camera) and 3 points which define a
// rectangle (a "window" into the scene)
type Camera struct {
	Focus      Vec3
	BottomLeft Vec3
	TopLeft    Vec3
	TopRight   Vec3
}

// ShootRay generates a ray coming out of the camera, going through the
// specified coordinates of the screen
func (c *Camera) ShootRay(x, y float64) *Ray {
	r := &Ray{}
	r.Start = c.Focus
	intersection := AddVectors(
		MinusVectors(&c.TopRight, &c.TopLeft).Scaled(x),
		MinusVectors(&c.BottomLeft, &c.TopLeft).Scaled(y),
	)
	r.Direction = *MinusVectors(intersection, &r.Start).Normalised()
	return r
}
