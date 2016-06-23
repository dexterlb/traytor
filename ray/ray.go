package ray

import (
	"fmt"
	"math"

	"github.com/DexterLB/traytor/maths"
)

// Ray is defined by its start, direction and depth which indicates how many materials it has passed through
type Ray struct {
	Start     maths.Vec3
	Direction maths.Vec3
	Depth     int
	Inverse   [3]float64
}

// New returns new ray
func New(start maths.Vec3, direction maths.Vec3, depth int) *Ray {
	return &Ray{Start: start, Direction: direction, Depth: depth}
}

//String returns the string representation of the ray
// in the form of "<start> -> <direction>"
func (r *Ray) String() string {
	return fmt.Sprintf("%s -> %s", &r.Start, &r.Direction)
}

//Init fills the Inverse field of ray
func (r *Ray) Init() {
	for i := 0; i < 3; i++ {
		r.Inverse[i] = 0
	}
	if math.Abs(r.Direction.X) > maths.Epsilon {
		r.Inverse[0] = 1.0 / r.Direction.X
	}
	if math.Abs(r.Direction.Y) > maths.Epsilon {
		r.Inverse[1] = 1.0 / r.Direction.Y
	}
	if math.Abs(r.Direction.Z) > maths.Epsilon {
		r.Inverse[2] = 1.0 / r.Direction.Z
	}
}

// Intersection represents a point on a surface struck by a ray
type Intersection struct {
	Point     *maths.Vec3
	Incoming  *Ray
	Material  int
	Distance  float64
	U, V      float64
	Normal    *maths.Vec3
	SurfaceOx *maths.Vec3
	SurfaceOy *maths.Vec3
}
