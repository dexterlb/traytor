package traytor

import (
	"math"
)

type BoundingBox struct {
	MaxVolume, MinVolume *Vec3
}

//NewBoundingBox creates a bounding box with no volume (min > max)
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinVolume: NewVec3(Inf, Inf, Inf),
		MaxVolume: NewVec3(-Inf, -Inf, -Inf),
	}
}

//AddPoint expands the volume of the box, if the point isn't already in the box
func (b *BoundingBox) AddPoint(point Vec3) {
	b.MinVolume.X = math.Min(b.MinVolume.X, point.X)
	b.MinVolume.Y = math.Min(b.MinVolume.Y, point.Y)
	b.MinVolume.Z = math.Min(b.MinVolume.Z, point.Z)

	b.MaxVolume.X = math.Max(b.MaxVolume.X, point.X)
	b.MaxVolume.Y = math.Max(b.MaxVolume.Y, point.Y)
	b.MaxVolume.Z = math.Max(b.MaxVolume.Z, point.Z)
}

//Inside checks if a point is inside the box
func (b *BoundingBox) Inside(point *Vec3) bool {
	return (Between(b.MinVolume.X, b.MaxVolume.X, point.X) &&
		Between(b.MinVolume.Y, b.MaxVolume.Y, point.Y) &&
		Between(b.MinVolume.Z, b.MaxVolume.Z, point.Z))
}

func (b *BoundingBox) IntersectAxis(ray *Ray, axis int) bool {
	//if the ray isn't pointing at the box there wouldn't be an intersection
	if (ray.Direction.GetDimension(axis) > 0 && ray.Start.GetDimension(axis) > b.MaxVolume.GetDimension(axis)) ||
		(ray.Direction.GetDimension(axis) < 0 && ray.Start.GetDimension(axis) < b.MinVolume.GetDimension(axis)) {
		return false
	}
	//or if the ray isn't moving in this direction at all
	if ray.Direction.GetDimension(axis) < Epsilon {
		return false
	}
	//we take the other two axes
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
	multiplier := ray.Direction.Inverse().GetDimension(axis)
	var intersectionX, intersectionY float64

	distance := (b.MinVolume.GetDimension(axis) - ray.Start.GetDimension(axis)) * multiplier
	if distance < 0 {
		return false
	}
	intersectionX = ray.Start.GetDimension(otherAxis1) + ray.Direction.GetDimension(otherAxis1)*distance
	if Between(b.MinVolume.GetDimension(otherAxis1), b.MaxVolume.GetDimension(otherAxis1), intersectionX) {
		intersectionY = ray.Start.GetDimension(otherAxis2) + ray.Direction.GetDimension(otherAxis2)*distance
		if Between(b.MinVolume.GetDimension(otherAxis2), b.MaxVolume.GetDimension(otherAxis2), intersectionY) {
			return true
		}
	}

	distance = (b.MaxVolume.GetDimension(axis) - ray.Start.GetDimension(axis)) * multiplier
	if distance < 0 {
		return false
	}
	intersectionX = ray.Start.GetDimension(otherAxis1) + ray.Direction.GetDimension(otherAxis1)*distance
	if Between(b.MinVolume.GetDimension(otherAxis1), b.MaxVolume.GetDimension(otherAxis1), intersectionX) {
		intersectionY = ray.Start.GetDimension(otherAxis2) + ray.Direction.GetDimension(otherAxis2)*distance
		if Between(b.MinVolume.GetDimension(otherAxis2), b.MaxVolume.GetDimension(otherAxis2), intersectionY) {
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
	var ray *Ray
	for rayStart := 0; rayStart < 3; rayStart++ {
		for rayEnd := rayStart + 1; rayEnd < 3; rayEnd++ {
			ray.Start = *triangle[rayStart]
			ray.Direction = *MinusVectors(triangle[rayEnd], triangle[rayStart])
			//Check if there's intersection between AB and the box (example)
			if b.Intersect(ray) {
				//we check whether there's intersection between BA and the Box too
				//to be sure the edge isn't outside the box
				//   _____
				//   |    | <----AB there is intersection, but BA----->
				//   |____|
				ray.Start = *triangle[rayEnd]
				ray.Direction = *MinusVectors(triangle[rayStart], triangle[rayEnd])
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
	var rayEnd *Vec3

	for edgeMask := 0; edgeMask < 7; edgeMask++ {
		for axis := uint(0); axis < 3; axis++ {
			if edgeMask&(1<<axis) != 0 {
				continue
			}

			if edgeMask&1 != 0 {
				ray.Start.X = b.MaxVolume.X
			} else {
				ray.Start.X = b.MinVolume.X
			}

			if edgeMask&2 != 0 {
				ray.Start.Y = b.MaxVolume.Y
			} else {
				ray.Start.Y = b.MinVolume.Y
			}

			if edgeMask&4 != 0 {
				ray.Start.Z = b.MaxVolume.Z
			} else {
				ray.Start.Z = b.MinVolume.Z
			}

			*rayEnd = ray.Start
			rayEnd.SetDimension(int(axis), b.MaxVolume.GetDimension(int(axis)))

			if (DotProduct(&ray.Start, ABxAC)-distance)*(DotProduct(&ray.Start, ABxAC)-distance) <= 0 {
				ray.Direction = *MinusVectors(rayEnd, &ray.Start)
				intersected, distance := IntersectTriangle(ray, A, B, C)
				if intersected && distance < 1.0000001 {
					return true
				}
			}
		}
	}

	return false
}

/*
line bool intersectTriangle(const Vector& A, const Vector& B, const Vector& C) const
		Vector AB = B - A;
		Vector AC = C - A;
		Vector ABcrossAC = AB ^ AC;
		double D = A * ABcrossAC;
		for (int mask = 0; mask < 7; mask++) {
			for (int j = 0; j < 3; j++) {
				if (mask & (1 << j)) continue;
				ray.start.set((mask & 1) ? vmax.x : vmin.x, (mask & 2) ? vmax.y : vmin.y, (mask & 4) ? vmax.z : vmin.z);
				Vector rayEnd = ray.start;
				rayEnd[j] = vmax[j];
				if (signOf(ray.start * ABcrossAC - D) != signOf(rayEnd * ABcrossAC - D)) {
					ray.dir = rayEnd - ray.start;
					ray.prepareForTracing();
					double gamma = 1.0000001;
					if (intersectTriangleFast(ray, A, B, C, gamma)) return true;
				}
			}
		}
		return false;
*/
