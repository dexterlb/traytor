package camera

import (
	"encoding/json"
	"fmt"

	"github.com/DexterLB/traytor/jsonutil"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// AnyCamera implements the Camera interface and is deserialiseable from json
type AnyCamera struct {
	Camera
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *AnyCamera) UnmarshalJSON(data []byte) error {
	cameraType, err := jsonutil.ObjectType(data)
	if err != nil {
		return err
	}

	switch cameraType {
	case "pinhole":
		camera := &PinholeCamera{}
		err := json.Unmarshal(data, &camera)
		if err != nil {
			return err
		}
		*c = AnyCamera{camera}
	default:
		return fmt.Errorf("Unknown camera type: '%s'", cameraType)
	}

	return nil
}

// Camera is a generic camera
type Camera interface {
	// ShootRay generates a ray which corresponds to the specified 2D coordinates
	// in the camera's viewframe
	ShootRay(x, y float64) *ray.Ray
}

// PinholeCamera has a focus
// (the location of the camera) and 3 points which define a
// rectangle (a "window" into the scene)
type PinholeCamera struct {
	Focus      maths.Vec3 `json:"focus"`
	BottomLeft maths.Vec3 `json:"bottom_left"`
	TopLeft    maths.Vec3 `json:"top_left"`
	TopRight   maths.Vec3 `json:"top_right"`
}

// ShootRay generates a ray coming out of the camera, going through the
// specified coordinates of the screen
func (c *PinholeCamera) ShootRay(x, y float64) *ray.Ray {
	r := &ray.Ray{}
	r.Start = c.Focus

	intersection := &maths.Vec3{}
	*intersection = c.TopLeft
	intersection.Add(maths.MinusVectors(&c.TopRight, &c.TopLeft).Scaled(x))
	intersection.Add(maths.MinusVectors(&c.BottomLeft, &c.TopLeft).Scaled(y))

	r.Direction = *maths.MinusVectors(intersection, &r.Start).Normalised()
	return r
}
