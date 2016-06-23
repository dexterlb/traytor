package mesh

import (
	"math"

	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// If the depth of the KDtree reaches MaxTreeDepth, the node becomes leaf
// whith node.Triangles the remaining triangles
const (
	MaxTreeDepth     = 10
	TrianglesPerLeaf = 20
)

// Vertex is a single vertex in a mesh
type Vertex struct {
	Normal      maths.Vec3 `json:"normal"`
	Coordinates maths.Vec3 `json:"coordinates"`
	UV          maths.Vec3 `json:"uv"`
}

// Triangle is a face with 3 vertices (indices in the vertex array)
type Triangle struct {
	Vertices      [3]int `json:"vertices"`
	Material      int    `json:"material"`
	AB, AC, ABxAC *maths.Vec3
	Normal        *maths.Vec3 `json:"normal"`
	surfaceOx     *maths.Vec3
	surfaceOy     *maths.Vec3
}

// Mesh is a triangle mesh
type Mesh struct {
	Vertices    []Vertex   `json:"vertices"`
	Faces       []Triangle `json:"faces"`
	tree        *KDtree
	BoundingBox *BoundingBox
}

// Init  of Mesh sets the Triangle indices in the Faces array, calculates the bounding box
// and KD tree, sets the surfaceOx and Oy and Cross Products of the sides of each triangle
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

		AB := maths.MinusVectors(B, A)
		AC := maths.MinusVectors(C, A)

		triangle.AB = AB
		triangle.AC = AC
		triangle.ABxAC = maths.CrossProduct(AB, AC)

		surfaceA := &m.Vertices[triangle.Vertices[0]].UV
		surfaceB := &m.Vertices[triangle.Vertices[1]].UV
		surfaceC := &m.Vertices[triangle.Vertices[2]].UV

		surfaceAB := maths.MinusVectors(surfaceB, surfaceA)
		surfaceAC := maths.MinusVectors(surfaceC, surfaceA)

		// Solve using Cramer:
		// |surfaceAB.X * px + surfaceAX.X *qx + 1 = 0
		// |surfaceAB.Y * px + surfaceAX.Y *qx = 0
		// and
		// surfaceAB.X * py + surfaceAX.X *qy = 0
		// surfaceAB.X * py + surfaceAX.X *qy + 1 = 0

		px, qx := maths.SolveEquation(surfaceAB, surfaceAC, maths.NewVec3(1, 0, 0))
		py, qy := maths.SolveEquation(surfaceAB, surfaceAC, maths.NewVec3(0, 1, 0))

		triangle.surfaceOx = maths.AddVectors(AB.Scaled(px), AC.Scaled(qx))
		triangle.surfaceOy = maths.AddVectors(AB.Scaled(py), AC.Scaled(qy))
	}
}

// SlowIntersect finds the intersection between a ray and the mesh
// and returns their intersection and the surface material.
// Returns nil and -1 if they don't intersect
// Has O(n) complexity.
func (m *Mesh) SlowIntersect(incoming *ray.Ray) *ray.Intersection {
	intersection := &ray.Intersection{}
	intersection.Distance = maths.Inf
	found := false
	for _, triangle := range m.Faces {
		if m.intersectTriangle(incoming, &triangle, intersection, nil) {
			found = true
		}
	}
	if !found {
		return nil
	}
	return intersection
}

// Intersect finds the intersection between a ray and the mesh
// and returns their intersection and the surface material.
// Returns nil and -1 if they don't intersect
// Has O(log(n)) amortised complexity.
func (m *Mesh) Intersect(incoming *ray.Ray) *ray.Intersection {
	incoming.Init()
	// There wouldn't be intersection if the incoming doesn't cross the bounding box
	if !m.BoundingBox.Intersect(incoming) {
		return nil
	}
	intersectionInfo := &ray.Intersection{Distance: maths.Inf}
	if m.IntersectKD(incoming, m.BoundingBox, m.tree, intersectionInfo) {
		return intersectionInfo
	}
	return nil
}

// IntersectTriangle find whether there's an intersection point between the ray and the triangle
// using barycentric coordinates and calculate the distance
func IntersectTriangle(ray *ray.Ray, A, B, C *maths.Vec3) (bool, float64) {
	AB := maths.MinusVectors(B, A)
	AC := maths.MinusVectors(C, A)
	reverseDirection := ray.Direction.Negative()
	distToA := maths.MinusVectors(&ray.Start, A)
	ABxAC := maths.CrossProduct(AB, AC)
	det := maths.DotProduct(ABxAC, reverseDirection)
	reverseDet := 1 / det
	if math.Abs(det) < maths.Epsilon {
		return false, maths.Inf
	}
	lambda2 := maths.MixedProduct(distToA, AC, reverseDirection) * reverseDet
	lambda3 := maths.MixedProduct(AB, distToA, reverseDirection) * reverseDet
	gamma := maths.DotProduct(ABxAC, distToA) * reverseDet
	if gamma < 0 {
		return false, maths.Inf
	}
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2+lambda3 > 1 {
		return false, maths.Inf
	}
	return true, gamma
}

// IntersectTriangle returns whether there's an intersection between the ray and the triangle,
// using barycentric coordinates and takes the point only if it's closer to the
// previously found intersection and the point is within the bounding box
func (m *Mesh) intersectTriangle(ray *ray.Ray, triangle *Triangle, intersection *ray.Intersection, boundingBox *BoundingBox) bool {
	// lambda2 * AB + lambda3 * AC - intersectDist*rayDir = distToA
	// If the triangle is ABC, this gives you A
	A := &m.Vertices[triangle.Vertices[0]].Coordinates
	distToA := maths.MinusVectors(&ray.Start, A)
	rayDir := ray.Direction
	ABxAC := triangle.ABxAC
	// We will find the barycentric coordinates using Cramer's formula, so we'll need the determinant
	// det is (AB^AC)*dir of the ray, but we're gonna use 1/det, so we find the recerse:
	det := -maths.DotProduct(ABxAC, &rayDir)
	if math.Abs(det) < maths.Epsilon {
		return false
	}
	reverseDet := 1 / det
	intersectDist := maths.DotProduct(ABxAC, distToA) * reverseDet
	if intersectDist < 0 || intersectDist > intersection.Distance {
		return false
	}
	// lambda2 = (dist^dir)*AC / det
	// lambda3 = -(dist^dir)*AB / det
	lambda2 := maths.MixedProduct(distToA, &rayDir, triangle.AC) * reverseDet
	lambda3 := -maths.MixedProduct(distToA, &rayDir, triangle.AB) * reverseDet
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2+lambda3 > 1 {
		return false
	}

	ip := maths.AddVectors(&ray.Start, (&rayDir).Scaled(intersectDist))

	// If we aren't inside the bounding box, there could be a closer intersection
	// within the bounding box
	if boundingBox != nil && !boundingBox.Inside(ip) {
		return false
	}
	intersection.Point = ip
	intersection.Distance = intersectDist
	if triangle.Normal != nil {
		intersection.Normal = triangle.Normal
	} else {
		// We solve intersection.normal = Anormal + AB normal * lambda2 + ACnormal * lambda 3
		Anormal := &m.Vertices[triangle.Vertices[0]].Normal
		Bnormal := &m.Vertices[triangle.Vertices[1]].Normal
		Cnormal := &m.Vertices[triangle.Vertices[2]].Normal
		ABxlambda2 := maths.MinusVectors(Bnormal, Anormal).Scaled(lambda2)
		ACxlambda3 := maths.MinusVectors(Cnormal, Anormal).Scaled(lambda3)
		intersection.Normal = maths.AddVectors(Anormal, maths.AddVectors(ABxlambda2, ACxlambda3))
	}
	uvA := &m.Vertices[triangle.Vertices[0]].UV
	uvB := &m.Vertices[triangle.Vertices[1]].UV
	uvC := &m.Vertices[triangle.Vertices[2]].UV

	// We solve intersection.uv = uvA + uvAB * lambda2 + uvAC * lambda 3
	uvABxlambda2 := maths.MinusVectors(uvB, uvA).Scaled(lambda2)
	uvACxlambda3 := maths.MinusVectors(uvC, uvA).Scaled(lambda3)
	uv := maths.AddVectors(uvA, maths.AddVectors(uvABxlambda2, uvACxlambda3))

	intersection.U = uv.X
	intersection.V = uv.Y

	intersection.SurfaceOx = triangle.surfaceOx
	intersection.SurfaceOy = triangle.surfaceOy

	intersection.Incoming = ray
	intersection.Material = triangle.Material
	return true
}

// GetBoundingBox returns the boundig box of the mesh, adding every vertex to the box
func (m *Mesh) GetBoundingBox() *BoundingBox {
	boundingBox := NewBoundingBox()
	for _, vertex := range m.Vertices {
		boundingBox.AddPoint(&vertex.Coordinates)
	}
	return boundingBox
}

// NewKDtree returns the KD tree for the mesh with MaxTreeDepth by slicing the bouindingBox
// and including the triangles in the bounding box (if it's in the middle of two bounding boxes, we include it in both)
func (m *Mesh) newKDtree(boundingBox *BoundingBox, trianglesIndices []int, depth int) *KDtree {
	if depth > MaxTreeDepth || len(trianglesIndices) < TrianglesPerLeaf {
		node := NewLeaf(trianglesIndices)
		return node
	}
	// We take the (axis + 2) % 3 to alternate between Ox, Oy and Oz on each turn
	// it's not the best decision in every case
	axis := (depth + 2) % 3
	leftLimit := boundingBox.MaxVolume[axis]
	righLimit := boundingBox.MinVolume[axis]

	median := (leftLimit + righLimit) / 2

	var leftTriangles, rightTriangles []int
	var A, B, C *maths.Vec3
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

// IntersectKD returns whether there's an intersection with the ray. The the current node is leaf
// we check each of its triangles and divide the bounding box and check for each child
func (m *Mesh) IntersectKD(ray *ray.Ray, boundingBox *BoundingBox, node *KDtree, intersectionInfo *ray.Intersection) bool {
	foundIntersection := false
	if node.Axis == maths.Leaf {
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
	}
	if firstBoundingBox.Intersect(ray) {
		return m.IntersectKD(ray, firstBoundingBox, firstNodeChild, intersectionInfo)
	}
	return m.IntersectKD(ray, secondBoundingBox, secondNodeChild, intersectionInfo)
}
