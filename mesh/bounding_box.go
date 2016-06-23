package mesh

import (
	"fmt"
	"math"

	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// BoundingBox represents a cuboid which surrounds a part of the mesh.
// The points in the cuboid are those (x, y, z), for which is true:
// MinVolume[0] < x < MaxVolume[0]
// MinVolume[1] < y < MaxVolume[1]
// MinVolume[2] < z , MaxVolume[2]
type BoundingBox struct {
	MaxVolume, MinVolume [3]float64
}

// NewBoundingBox creates a bounding box with no volume (min > max)
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinVolume: [3]float64{maths.Inf, maths.Inf, maths.Inf},
		MaxVolume: [3]float64{-maths.Inf, -maths.Inf, -maths.Inf},
	}
}

// AddPoint expands the volume of the box, if the point isn't already in the box
func (b *BoundingBox) AddPoint(point *maths.Vec3) {
	b.MinVolume[0] = math.Min(b.MinVolume[0], point.X)
	b.MinVolume[1] = math.Min(b.MinVolume[1], point.Y)
	b.MinVolume[2] = math.Min(b.MinVolume[2], point.Z)

	b.MaxVolume[0] = math.Max(b.MaxVolume[0], point.X)
	b.MaxVolume[1] = math.Max(b.MaxVolume[1], point.Y)
	b.MaxVolume[2] = math.Max(b.MaxVolume[2], point.Z)
}

// Inside checks if a point is inside the box
func (b *BoundingBox) Inside(point *maths.Vec3) bool {
	return (maths.Between(b.MinVolume[0], b.MaxVolume[0], point.X) &&
		maths.Between(b.MinVolume[1], b.MaxVolume[1], point.Y) &&
		maths.Between(b.MinVolume[2], b.MaxVolume[2], point.Z))
}

// otherAxes returns the two other axis of the given and we can get (0, 1, 2), (1, 0, 2), (2, 0, 1)
// if axis is the first of the tree, otherAxis1 is the second, and otherAxis2 the third
func otherAxes(axis int) (int, int) {
	var otherAxis1, otherAxis2 int
	if axis == 0 {
		otherAxis1 = 1
	} else {
		otherAxis1 = 0
	}
	if axis == 2 {
		otherAxis2 = 1
	} else {
		otherAxis2 = 2
	}
	return otherAxis1, otherAxis2
}

// IntersectAxis checks whether there's interaction between the ray and the box.
func (b *BoundingBox) IntersectAxis(ray *ray.Ray, axis int) bool {
	directions := [3]float64{ray.Direction.X, ray.Direction.Y, ray.Direction.Z}
	start := [3]float64{ray.Start.X, ray.Start.Y, ray.Start.Z}

	// if the ray isn't pointing at the box there wouldn't be an intersection
	if (directions[axis] > 0 && start[axis] > b.MaxVolume[axis]) ||
		(directions[axis] < 0 && start[axis] < b.MinVolume[axis]) {
		return false
	}
	// or if the ray isn't moving in this direction at all
	if math.Abs(directions[axis]) < maths.Epsilon {
		return false
	}
	// we take the other two axes
	otherAxis1, otherAxis2 := otherAxes(axis)

	multiplier := ray.Inverse[axis]
	var intersectionX, intersectionY float64

	distance := (b.MinVolume[axis] - start[axis]) * multiplier
	if distance < 0 {
		return false
	}

	intersectionX = start[otherAxis1] + directions[otherAxis1]*distance
	if maths.Between(b.MinVolume[otherAxis1], b.MaxVolume[otherAxis1], intersectionX) {
		intersectionY = start[otherAxis2] + directions[otherAxis2]*distance
		if maths.Between(b.MinVolume[otherAxis2], b.MaxVolume[otherAxis2], intersectionY) {
			return true
		}
	}

	distance = (b.MaxVolume[axis] - start[axis]) * multiplier
	if distance < 0 {
		return false
	}
	intersectionX = start[otherAxis1] + directions[otherAxis1]*distance
	if maths.Between(b.MinVolume[otherAxis1], b.MaxVolume[otherAxis1], intersectionX) {
		intersectionY = start[otherAxis2] + directions[otherAxis2]*distance
		if maths.Between(b.MinVolume[otherAxis2], b.MaxVolume[otherAxis2], intersectionY) {
			return true
		}
	}
	return false
}

// Intersect check if a ray intersects the bounding box
func (b *BoundingBox) Intersect(ray *ray.Ray) bool {
	if b.Inside(&ray.Start) {
		return true
	}
	return (b.IntersectAxis(ray, maths.Ox) || b.IntersectAxis(ray, maths.Oy) || b.IntersectAxis(ray, maths.Oz))
}

// IntersectTriangle checks if the bounding box intersects with a triangle
// 1) To have a vertex in the box
// 2) The edge of the triangle intersects with the box
// 3) The middle of the triangle to be inside the box, while the vertices aren't
func (b *BoundingBox) IntersectTriangle(A, B, C *maths.Vec3) bool {
	if b.Inside(A) || b.Inside(B) || b.Inside(C) {
		return true
	}
	// we construct the ray from A->B, A->C, etc and check whether it intersects with the box
	triangle := [3]*maths.Vec3{A, B, C}
	ray := &ray.Ray{}
	for rayStart := 0; rayStart < 3; rayStart++ {
		for rayEnd := rayStart + 1; rayEnd < 3; rayEnd++ {
			ray.Start = *(triangle[rayStart])
			ray.Direction = *maths.MinusVectors(triangle[rayEnd], triangle[rayStart])
			ray.Init()
			// Check if there's intersection between AB and the box (example)
			if b.Intersect(ray) {
				// we check whether there's intersection between BA and the Box too
				// to be sure the edge isn't outside the box
				//    _____
				//    |    | <----AB there is intersection, but BA----->
				//    |____|
				ray.Start = *triangle[rayEnd]
				ray.Direction = *maths.MinusVectors(triangle[rayStart], triangle[rayEnd])
				ray.Init()
				if b.Intersect(ray) {
					return true
				}
			}
		}
	}
	// we have to check if the inside of the triangle intersects with the box
	AB := maths.MinusVectors(B, A)
	AC := maths.MinusVectors(C, A)
	ABxAC := maths.CrossProduct(AB, AC)
	distance := maths.DotProduct(A, ABxAC)
	rayEnd := &maths.Vec3{}

	for edgeMask := 0; edgeMask < 7; edgeMask++ {
		for axis := uint(0); axis < 3; axis++ {
			if edgeMask&(1<<axis) != 0 {
				continue
			}

			if edgeMask&1 != 0 {
				ray.Start.X = b.MaxVolume[0]
			} else {
				ray.Start.X = b.MinVolume[0]
			}

			if edgeMask&2 != 0 {
				ray.Start.Y = b.MaxVolume[1]
			} else {
				ray.Start.Y = b.MinVolume[1]
			}

			if edgeMask&4 != 0 {
				ray.Start.Z = b.MaxVolume[2]
			} else {
				ray.Start.Z = b.MinVolume[2]
			}

			*rayEnd = ray.Start
			rayEnd.SetDimension(int(axis), b.MaxVolume[axis])

			if (maths.DotProduct(&ray.Start, ABxAC)-distance)*(maths.DotProduct(rayEnd, ABxAC)-distance) <= 0 {
				ray.Direction = *maths.MinusVectors(rayEnd, &ray.Start)
				ray.Init()
				intersected, distance := IntersectTriangle(ray, A, B, C)
				if intersected && distance < 1.0000001 {
					return true
				}
			}
		}
	}

	return false
}

// Split returns two new bounding boxes which are the result of spliting the original on the given axis and median
func (b *BoundingBox) Split(axis int, median float64) (*BoundingBox, *BoundingBox) {
	left := &BoundingBox{}
	*left = *b
	right := &BoundingBox{}
	*right = *b
	left.MaxVolume[axis] = median
	right.MinVolume[axis] = median
	return left, right
}

// IntersectWall checks if a ray intersects a wall inside the bounding box
// (the wall is defined by the axis and median, as with Split)
func (b *BoundingBox) IntersectWall(axis int, median float64, ray *ray.Ray) bool {
	directions := [3]float64{ray.Direction.X, ray.Direction.Y, ray.Direction.Z}
	start := [3]float64{ray.Start.X, ray.Start.Y, ray.Start.Z}
	if math.Abs(directions[axis]) < maths.Epsilon {
		return (math.Abs(start[axis]-median) < maths.Epsilon)
	}

	otherAxis1, otherAxis2 := otherAxes(axis)
	distance := median - start[axis]
	directionInAxis := ray.Inverse[axis]

	if (distance * directionInAxis) < 0 {
		return false
	}

	fac := distance * directionInAxis
	distanceOnAxis1 := start[otherAxis1] +
		directions[otherAxis1]*fac
	if b.MinVolume[otherAxis1] <= distanceOnAxis1 &&
		distanceOnAxis1 <= b.MaxVolume[otherAxis1] {

		distanceOnAxis2 := start[otherAxis2] +
			directions[otherAxis2]*fac
		return b.MinVolume[otherAxis2] <= distanceOnAxis2 &&
			distanceOnAxis2 <= b.MaxVolume[otherAxis2]
	}
	return false
}

// String returns the string representation of the boundingBox in the form of bbox[min: _, max: _]
func (b *BoundingBox) String() string {
	return fmt.Sprintf("bbox[min: %s, max: %s]", maths.NewVec3Array(b.MinVolume), maths.NewVec3Array(b.MaxVolume))
}
