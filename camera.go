package traytor

// Camera is a generic camera
type Camera interface {
	// ShootRay generates a ray which corresponds to the specified 2D coordinates
	// in the camera's viewframe
	ShootRay(x, y float64) *Ray
}

// PinholeCamera has a focus
// (the location of the camera) and 3 points which define a
// rectangle (a "window" into the scene)
type PinholeCamera struct {
	Focus      Vec3 `json:"focus"`
	BottomLeft Vec3 `json:"bottom_left"`
	TopLeft    Vec3 `json:"top_left"`
	TopRight   Vec3 `json:"top_right"`
}

// ShootRay generates a ray coming out of the camera, going through the
// specified coordinates of the screen
func (c *PinholeCamera) ShootRay(x, y float64) *Ray {
	r := &Ray{}
	r.Start = c.Focus

	intersection := &Vec3{}
	*intersection = c.TopLeft
	intersection.Add(MinusVectors(&c.TopRight, &c.TopLeft).Scaled(x))
	intersection.Add(MinusVectors(&c.BottomLeft, &c.TopLeft).Scaled(y))

	r.Direction = *MinusVectors(intersection, &r.Start).Normalised()
	return r
}
