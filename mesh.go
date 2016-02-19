package traytor

import (
	//"log"
	"math"
)

const (
	MaxTreeDepth = 10
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
	Vertices    []Vertex   `json:"vertices"`
	Faces       []Triangle `json:"faces"`
	tree        *KDtree
	BoundingBox *BoundingBox
}

func (m *Mesh) Init() {
	allIndices := make([]int, len(m.Faces))
	for i := range allIndices {
		allIndices[i] = i
	}

	m.BoundingBox = m.GetBoundingBox()
	m.tree = m.newKDtree(m.BoundingBox, allIndices, 0)
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
func (m *Mesh) SlowIntersect(ray *Ray) *Intersection {
	intersection := &Intersection{}
	intersection.Distance = Inf
	found := false
	for _, triangle := range m.Faces {
		if m.intersectTriangle(ray, &triangle, intersection, nil) {
			found = true
		}
	}
	if !found {
		return nil
	}
	return intersection
}

func (m *Mesh) Intersect(ray *Ray) *Intersection {
	if !m.BoundingBox.Intersect(ray) {
		return nil
	}
	intersectionInfo := &Intersection{Distance: Inf}
	if m.IntersectKD(ray, m.BoundingBox, m.tree, intersectionInfo) {
		return intersectionInfo
	}
	return nil
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

func (m *Mesh) intersectTriangle(ray *Ray, triangle *Triangle, intersection *Intersection, boundingBox *BoundingBox) bool {
	//lambda2(B - A) + lambda3(C - A) - intersectDist*rayDir = distToA
	//If the triangle is ABC, this gives you A
	A := &m.Vertices[triangle.Vertices[0]].Coordinates
	distToA := MinusVectors(&ray.Start, A)
	rayDir := ray.Direction
	ABxAC := triangle.ABxAC
	//We will find the barycentric coordinates using Cramer's formula, so we'll need the determinant
	//det is (AB^AC)*dir of the ray, but we're gonna use 1/det, so we find the recerse:
	det := -DotProduct(ABxAC, &rayDir)
	if math.Abs(det) < Epsilon {
		return false
	}
	reverseDet := 1 / det
	intersectDist := DotProduct(ABxAC, distToA) * reverseDet
	if intersectDist < 0 || intersectDist > intersection.Distance {
		return false
	}
	//lambda2 = (dist^dir)*AC / det
	//lambda3 = -(dist^dir)*AB / det
	lambda2 := MixedProduct(distToA, &rayDir, triangle.AC) * reverseDet
	lambda3 := -MixedProduct(distToA, &rayDir, triangle.AB) * reverseDet
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2+lambda3 > 1 {
		return false
	}

	ip := AddVectors(&ray.Start, (&rayDir).Scaled(intersectDist))

	if boundingBox != nil && !boundingBox.Inside(ip) {
		return false
	}
	intersection.Point = ip
	intersection.Distance = intersectDist
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
	intersection.Material = triangle.Material
	return true
}

func (m *Mesh) GetBoundingBox() *BoundingBox {
	boundingBox := NewBoundingBox()
	for _, vertex := range m.Vertices {
		boundingBox.AddPoint(&vertex.Coordinates)
	}
	return boundingBox
}

func (m *Mesh) newKDtree(boundingBox *BoundingBox, trianglesIndices []int, depth int) *KDtree {
	if depth > MaxTreeDepth || len(trianglesIndices) < 2 {
		node := NewLeaf(trianglesIndices)
		return node
	}
	axis := (depth + 2) % 3
	leftLimit := boundingBox.MaxVolume.GetDimension(axis)
	righLimit := boundingBox.MinVolume.GetDimension(axis)

	median := (leftLimit + righLimit) / 2

	var leftTriangles, rightTriangles []int
	var A, B, C *Vec3
	leftBoundingBox, rightBoundingBox := boundingBox.Split(axis, median)
	for _, index := range trianglesIndices {
		A = &m.Vertices[m.Faces[index].Vertices[0]].Coordinates
		B = &m.Vertices[m.Faces[index].Vertices[1]].Coordinates
		C = &m.Vertices[m.Faces[index].Vertices[2]].Coordinates

		if leftBoundingBox.IntersectTriangle(A, B, C) {
			leftTriangles = append(leftTriangles, index)
		}

		if rightBoundingBox.IntersectTriangle(A, B, C) {
			rightTriangles = append(rightTriangles, index)
		}
	}
	node := NewNode(median, axis)
	leftChild := m.newKDtree(leftBoundingBox, leftTriangles, depth+1)
	rightChild := m.newKDtree(rightBoundingBox, rightTriangles, depth+1)
	node.Children[0] = leftChild
	node.Children[1] = rightChild
	return node
}

func (m *Mesh) IntersectKD(ray *Ray, boundingBox *BoundingBox, node *KDtree, intersectionInfo *Intersection) bool {
	foundIntersection := false
	if node.Axis == Leaf {
		for _, triangle := range node.Triangles {
			if m.intersectTriangle(ray, &m.Faces[triangle], intersectionInfo, boundingBox) {
				foundIntersection = true
			}
		}
		return foundIntersection
	}

	leftBoundingBoxChild, rightBoundingBoxChild := boundingBox.Split(node.Axis, node.Median)

	var firstBoundingBox, secondBoundingBox *BoundingBox
	var firstNodeChild, secondNodeChild *KDtree
	if ray.Start.GetDimension(node.Axis) <= node.Median {
		firstBoundingBox = leftBoundingBoxChild
		secondBoundingBox = rightBoundingBoxChild
		firstNodeChild = node.Children[0]
		secondNodeChild = node.Children[1]
	} else {
		firstBoundingBox = rightBoundingBoxChild
		secondBoundingBox = leftBoundingBoxChild
		firstNodeChild = node.Children[1]
		secondNodeChild = node.Children[0]
	}

	if boundingBox.IntersectWall(node.Axis, node.Median, ray) {
		if m.IntersectKD(ray, firstBoundingBox, firstNodeChild, intersectionInfo) {
			return true
		}
		return m.IntersectKD(ray, secondBoundingBox, secondNodeChild, intersectionInfo)
	} else {
		if firstBoundingBox.Intersect(ray) {
			return m.IntersectKD(ray, firstBoundingBox, firstNodeChild, intersectionInfo)
		} else {
			return m.IntersectKD(ray, secondBoundingBox, secondNodeChild, intersectionInfo)
		}
	}
	return false
}
