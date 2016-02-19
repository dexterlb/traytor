package traytor

import (
	"fmt"
	"math"
)

type BoundingBox struct {
	MaxVolume, MinVolume [3]float64
}

//NewBoundingBox creates a bounding box with no volume (min > max)
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinVolume: [3]float64{Inf, Inf, Inf},
		MaxVolume: [3]float64{-Inf, -Inf, -Inf},
	}
}

//AddPoint expands the volume of the box, if the point isn't already in the box
func (b *BoundingBox) AddPoint(point *Vec3) {
	b.MinVolume[0] = math.Min(b.MinVolume[0], point.X)
	b.MinVolume[1] = math.Min(b.MinVolume[1], point.Y)
	b.MinVolume[2] = math.Min(b.MinVolume[2], point.Z)

	b.MaxVolume[0] = math.Max(b.MaxVolume[0], point.X)
	b.MaxVolume[1] = math.Max(b.MaxVolume[1], point.Y)
	b.MaxVolume[2] = math.Max(b.MaxVolume[2], point.Z)
}

//Inside checks if a point is inside the box
func (b *BoundingBox) Inside(point *Vec3) bool {
	return (Between(b.MinVolume[0], b.MaxVolume[0], point.X) &&
		Between(b.MinVolume[1], b.MaxVolume[1], point.Y) &&
		Between(b.MinVolume[2], b.MaxVolume[2], point.Z))
}

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

func (b *BoundingBox) IntersectAxis(ray *Ray, axis int) bool {
	directions := [3]float64{ray.Direction.X, ray.Direction.Y, ray.Direction.Z}
	start := [3]float64{ray.Start.X, ray.Start.Y, ray.Start.Z}

	//if the ray isn't pointing at the box there wouldn't be an intersection
	if (directions[axis] > 0 && start[axis] > b.MaxVolume[axis]) ||
		(directions[axis] < 0 && start[axis] < b.MinVolume[axis]) {
		return false
	}
	//or if the ray isn't moving in this direction at all
	if math.Abs(directions[axis]) < Epsilon {
		return false
	}
	//we take the other two axes
	otherAxis1, otherAxis2 := otherAxes(axis)

	multiplier := ray.Inverse[axis]
	var intersectionX, intersectionY float64

	distance := (b.MinVolume[axis] - start[axis]) * multiplier
	if distance < 0 {
		return false
	}
	intersectionX = start[otherAxis1] + directions[otherAxis1]*distance
	if Between(b.MinVolume[otherAxis1], b.MaxVolume[otherAxis1], intersectionX) {
		intersectionY = start[otherAxis2] + directions[otherAxis2]*distance
		if Between(b.MinVolume[otherAxis2], b.MaxVolume[otherAxis2], intersectionY) {
			return true
		}
	}

	distance = (b.MaxVolume[axis] - start[axis]) * multiplier
	if distance < 0 {
		return false
	}
	intersectionX = start[otherAxis1] + directions[otherAxis1]*distance
	if Between(b.MinVolume[otherAxis1], b.MaxVolume[otherAxis1], intersectionX) {
		intersectionY = start[otherAxis2] + directions[otherAxis2]*distance
		if Between(b.MinVolume[otherAxis2], b.MaxVolume[otherAxis2], intersectionY) {
			return true
		}
	}
	return false
}

//Intersect check if a ray intersects the bounding box
func (b *BoundingBox) Intersect(ray *Ray) bool {
	if b.Inside(&ray.Start) {
		return true
	}
	return (b.IntersectAxis(ray, Ox) || b.IntersectAxis(ray, Oy) || b.IntersectAxis(ray, Oz))
}

//IntersectTriangle checks if the bounding box intersects with a triangle
//1) To have a vertex in the box
//2) The edge of the triangle intersects with the box
//3) The middle of the triangle to be inside the box, while the vertices aren't
func (b *BoundingBox) IntersectTriangle(A, B, C *Vec3) bool {
	if b.Inside(A) || b.Inside(B) || b.Inside(C) {
		return true
	}
	//we construct the ray from A->B, A->C, etc and check whether it intersects with the box
	triangle := [3]*Vec3{A, B, C}
	ray := &Ray{}
	for rayStart := 0; rayStart < 3; rayStart++ {
		for rayEnd := rayStart + 1; rayEnd < 3; rayEnd++ {
			ray.Start = *(triangle[rayStart])
			ray.Direction = *MinusVectors(triangle[rayEnd], triangle[rayStart])
			ray.Init()
			//Check if there's intersection between AB and the box (example)
			if b.Intersect(ray) {
				//we check whether there's intersection between BA and the Box too
				//to be sure the edge isn't outside the box
				//   _____
				//   |    | <----AB there is intersection, but BA----->
				//   |____|
				ray.Start = *triangle[rayEnd]
				ray.Direction = *MinusVectors(triangle[rayStart], triangle[rayEnd])
				ray.Init()
				if b.Intersect(ray) {
					return true
				}
			}
		}
	}
	//we have to check if the inside of the triangle intersects with the box
	AB := MinusVectors(B, A)
	AC := MinusVectors(C, A)
	ABxAC := CrossProduct(AB, AC)
	distance := DotProduct(A, ABxAC)
	rayEnd := &Vec3{}

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

			if (DotProduct(&ray.Start, ABxAC)-distance)*(DotProduct(rayEnd, ABxAC)-distance) <= 0 {
				ray.Direction = *MinusVectors(rayEnd, &ray.Start)
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
func (b *BoundingBox) IntersectWall(axis int, median float64, ray *Ray) bool {
	directions := [3]float64{ray.Direction.X, ray.Direction.Y, ray.Direction.Z}
	start := [3]float64{ray.Start.X, ray.Start.Y, ray.Start.Z}
	if math.Abs(directions[axis]) < Epsilon {
		return (math.Abs(start[axis]-median) < Epsilon)
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

func (b *BoundingBox) String() string {
	return fmt.Sprintf("bbox[min: %s, max: %s]", NewVec3Array(b.MinVolume), NewVec3Array(b.MaxVolume))
}
