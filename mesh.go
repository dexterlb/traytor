package traytor

import (
	"math"
)

// Vertex is a single vertex in a mesh
type Vertex struct {
	Normal      Vec3 `json:"normal"`
	Coordinates Vec3 `json:"coordinates"`
	UV          Vec3 `json:"uv"`
}

// Triangle is a face with 3 vertices (indices in the vertex array)
type Triangle struct {
	Vertices      [3]int `json:"vertices"`
	Material      int    `json:"material"`
	AB, AC, ABxAC *Vec3
	Normal        *Vec3 `json:"normal"`
	surfaceOx     *Vec3
	surfaceOy     *Vec3
}

// Mesh is a triangle mesh
type Mesh struct {
	Vertices []Vertex   `json:"vertices"`
	Faces    []Triangle `json:"faces"`
}

func (m *Mesh) Init() {
	for i := range m.Faces {
		triangle := &m.Faces[i]

		A := &m.Vertices[triangle.Vertices[0]].Coordinates
		B := &m.Vertices[triangle.Vertices[1]].Coordinates
		C := &m.Vertices[triangle.Vertices[2]].Coordinates

		AB := MinusVectors(B, A)
		AC := MinusVectors(C, A)

		triangle.AB = AB
		triangle.AC = AC
		triangle.ABxAC = CrossProduct(AB, AC)

		surfaceA := &m.Vertices[triangle.Vertices[0]].UV
		surfaceB := &m.Vertices[triangle.Vertices[1]].UV
		surfaceC := &m.Vertices[triangle.Vertices[2]].UV

		surfaceAB := MinusVectors(surfaceB, surfaceA)
		surfaceAC := MinusVectors(surfaceC, surfaceA)

		px, qx := SolveEquation(surfaceAB, surfaceAC, NewVec3(1, 0, 0))
		py, qy := SolveEquation(surfaceAB, surfaceAC, NewVec3(0, 1, 0))

		triangle.surfaceOx = AddVectors(AB.Scaled(px), AC.Scaled(qx))
		triangle.surfaceOy = AddVectors(AB.Scaled(py), AC.Scaled(qy))
	}
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
	if maxDistance >= Inf-Epsilon {
		return nil, -1
	}
	return intersection, material
}

//Intersect triangle looks very similar to mesh.IntersectTriangle
func IntersectTriangle(ray *Ray, A, B, C *Vec3) (bool, float64) {
	AB := MinusVectors(B, A)
	AC := MinusVectors(C, A)
	reverseDirection := ray.Direction.Negative()
	distToA := MinusVectors(&ray.Start, A)
	ABxAC := CrossProduct(AB, AC)
	det := DotProduct(ABxAC, reverseDirection)
	reverseDet := 1 / det
	if math.Abs(det) < Epsilon {
		return false, Inf
	}
	lambda2 := MixedProduct(distToA, AC, reverseDirection) * reverseDet
	lambda3 := MixedProduct(AB, distToA, reverseDirection) * reverseDet
	gamma := DotProduct(ABxAC, distToA) * reverseDet
	if gamma < 0 {
		return false, Inf
	}
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2+lambda3 > 1 {
		return false, Inf
	}
	return true, gamma
}

func (m *Mesh) intersectTriangle(ray *Ray, triangle *Triangle) (*Intersection, float64) {
	//lambda2(B - A) + lambda3(C - A) - intersectDist*rayDir = distToA
	intersection := &Intersection{}
	//If the triangle is ABC, this gives you A
	A := &m.Vertices[triangle.Vertices[0]].Coordinates
	distToA := MinusVectors(&ray.Start, A)
	rayDir := ray.Direction
	ABxAC := triangle.ABxAC
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
	lambda2 := MixedProduct(distToA, &rayDir, triangle.AC) * reverseDet
	lambda3 := -MixedProduct(distToA, &rayDir, triangle.AB) * reverseDet
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

	intersection.SurfaceOx = triangle.surfaceOx
	intersection.SurfaceOy = triangle.surfaceOy

	intersection.Incoming = ray
	return intersection, intersection.Distance
}

func (m *Mesh) GetBoundingBox() *BoundingBox {
	boundingBox := NewBoundingBox()
	for _, vertex := range m.Vertices {
		boundingBox.AddPoint(vertex.Coordinates)
	}
	return boundingBox
}

func (m *Mesh) KDTree(boundingBox *BoundingBox, trianglesIndices []int, depth int) *KDtree {
	node = &KDtree{}
	if depth > MaxDepth {
		node = NewLeaf(triangles)
		return node
	}
	axis := depth % 3
	leftLimit := boundingBox.MaxVolume.GetDimension(axis)
	righLimit := boundingBox.MinVolume.GetDimension(axis)
	median := (leftLimit + righLimit) / 2
	var leftTriangles, rightTriangles []int
	var A, B, C *Vec3
	leftBoundingBox, rightBoundingBox := boundingBox.Split(axis, median)
	for _, index := range trianglesIndices {
		A = m.Vertices[m.Faces[index].Vertices[0]].Coordinates
		B = m.Vertices[m.Faces[index].Vertices[1]].Coordinates
		C = m.Vertices[m.Faces[index].Vertices[1]].Coordinates

		if leftBoundingBox.IntersectTriangle(A, B, C) {
			leftTriangles = append(leftTriangles, index)
		}

		if rightBoundingBox.IntersectTriangle(A, B, C) {
			appendTriangles = append(rightTriangles, index)
		}
	}
	node = NewNode(median, axis)
	leftChild := m.KDTree(leftBoundingBox, leftTriangles, depth+1)
	rightChild := m.KDTree(rightBoundingBox, rightTriangles, depth+1)
	node.Children[0] = leftChild
	node.Children[1] = rightChild
	return node
}
