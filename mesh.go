package traytor

// Vertex is a single vertex in a mesh
type Vertex struct {
	Normal      Vec3
	Coordinates Vec3
	UV          Vec3
}

// Triangle is a face with 3 vertices (indices in the vertex array)
type Triangle struct {
	Points [3]int
}

// Mesh is a triangle mesh
type Mesh struct {
	Vertices []Vertex
	Faces    []Triangle
}
