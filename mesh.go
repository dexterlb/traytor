package traytor

import (
	"fmt"
	"math"
	"os"
)

// Vertex is a single vertex in a mesh
type Vertex struct {
	Normal      Vec3 `json:"normal"`
	Coordinates Vec3 `json:"coordinates"`
	UV          Vec3 `json:"uv"`
}

// Triangle is a face with 3 vertices (indices in the vertex array)
type Triangle struct {
	Vertices [3]int `json:"vertices"`
	Material int    `json:"material"`
	Normal   *Vec3  `json:"normal"`
}

// Mesh is a triangle mesh
type Mesh struct {
	Vertices []Vertex   `json:"vertices"`
	Faces    []Triangle `json:"faces"`
}

// Intersect finds the intersection between a ray and the mesh
// and returns their intersection and the surface material.
// Returns nil and -1 if they don't intersect
func (m *Mesh) Intersect(ray *Ray) (*Intersection, int) {
	maxDistance := Inf
	intersection := &Intersection{}
	tempIntersection := &Intersection{}
	material := -1
	var distance float64
	for _, triangle := range m.Faces {
		tempIntersection, distance = m.intersectTriangle(ray, &triangle)
		if distance < maxDistance {
			maxDistance = distance
			intersection = tempIntersection
			material = triangle.Material
		}
	}
	if distance >= Inf-Epsilon {
		return nil, -1
	}
	return intersection, material
}

func (m *Mesh) intersectTriangle(ray *Ray, triangle *Triangle) (*Intersection, float64) {
	//lambda2(B - A) + lambda3(C - A) - intersectDist*rayDir = distToA
	intersection := &Intersection{}
	//If the triangle is ABC, this gives you A
	A := &m.Vertices[triangle.Vertices[0]].Coordinates
	B := &m.Vertices[triangle.Vertices[1]].Coordinates
	C := &m.Vertices[triangle.Vertices[2]].Coordinates

	distToA := MinusVectors(&ray.Start, A)
	rayDir := ray.Direction
	ABxAC := CrossProduct(MinusVectors(B, A), MinusVectors(C, A))
	//We will find the barycentric coordinates using Cramer's formula, so we'll need the determinant
	//det is (AB^AC)*dir of the ray, but we're gonna use 1/det, so we find the recerse:
	det := -DotProduct(ABxAC, &rayDir)
	if math.Abs(det) < Epsilon {
		return nil, Inf
	}
	reverseDet := 1 / det
	intersectDist := DotProduct(ABxAC, distToA) * reverseDet
	if intersectDist < 0 {
		return nil, Inf
	}
	//lambda2 = (dist^dir)*AC / det
	//lambda3 = -(dist^dir)*AB / det
	fmt.Fprintf(os.Stderr, "mixed: %s\n", distToA)
	lambda2 := MixedProduct(distToA, &rayDir, MinusVectors(C, A)) * reverseDet
	lambda3 := -MixedProduct(distToA, &rayDir, MinusVectors(B, A)) * reverseDet
	fmt.Fprintf(os.Stderr, "lambda2: %.3g, lambda3: %.3g", lambda2, lambda3)
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2+lambda3 > 1 {
		return nil, Inf
	}

	intersection.Distance = intersectDist
	intersection.Point = AddVectors(&ray.Start, (&rayDir).Scaled(intersectDist))
	if triangle.Normal != nil {
		intersection.Normal = triangle.Normal
	} else {
		Anormal := &m.Vertices[triangle.Vertices[0]].Normal
		Bnormal := &m.Vertices[triangle.Vertices[1]].Normal
		Cnormal := &m.Vertices[triangle.Vertices[2]].Normal
		ABxlambda2 := MinusVectors(Bnormal, Anormal).Scaled(lambda2)
		ACxlambda3 := MinusVectors(Cnormal, Anormal).Scaled(lambda3)
		intersection.Normal = AddVectors(Anormal, AddVectors(ABxlambda2, ACxlambda3))
	}
	uvA := &m.Vertices[triangle.Vertices[0]].UV
	uvB := &m.Vertices[triangle.Vertices[1]].UV
	uvC := &m.Vertices[triangle.Vertices[2]].UV

	uvABxlambda2 := MinusVectors(uvB, uvA).Scaled(lambda2)
	uvACxlambda3 := MinusVectors(uvC, uvA).Scaled(lambda3)
	uv := AddVectors(uvA, AddVectors(uvABxlambda2, uvACxlambda3))
	intersection.U = uv.X
	intersection.V = uv.Y

	intersection.Incoming = ray
	return intersection, intersection.Distance
}
