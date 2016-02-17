package traytor

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
	Normal   Vec3   `json:"normal"`
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
	return nil, -1
}

func (m *Mesh) IntersectionTriangle(ray *Ray, triangle *Triangle) *Intersection {
	//lambda2(B - A) + lambda3(C - A) - intersectDist*rayDir = distToA
	if DotProduct(&ray.Direction, &triangle.Normal) > 0 {
		return nil
	}
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
	det := DotProduct(ABxAC, &rayDir)
	if det < Epsilon {
		return nil
	}
	reverseDet := 1 / det
	intersectDist := DotProduct(ABxAC, distToA) * reverseDet
	intersection.Distance = intersectDist
	return intersection
}

/*
bool Mesh::intersectTriangle(const RRay& ray, const Triangle& t, IntersectionInfo& info)
{
	if (backfaceCulling && dot(ray.dir, t.gnormal) > 0) return false;
	Vector A = vertices[t.v[0]];

	Vector H = ray.start - A;
	Vector D = ray.dir;

	double Dcr = - (t.ABcrossAC * D);

	if (fabs(Dcr) < 1e-12) return false;

	double rDcr = 1 / Dcr;
	double gamma = (t.ABcrossAC * H) * rDcr;
	if (gamma < 0 || gamma > info.distance) return false;

	Vector HcrossD = H^D;
	double lambda2 = (HcrossD * t.AC) * rDcr;
	if (lambda2 < 0 || lambda2 > 1) return false;

	double lambda3 = -(t.AB * HcrossD) * rDcr;
	if (lambda3 < 0 || lambda3 > 1) return false;

	if (lambda2 + lambda3 > 1) return false;

	info.distance = gamma;
	info.ip = ray.start + ray.dir * gamma;
	if (!faceted) {
		Vector nA = normals[t.n[0]];
		Vector nB = normals[t.n[1]];
		Vector nC = normals[t.n[2]];

		info.normal = nA + (nB - nA) * lambda2 + (nC - nA) * lambda3;
		info.normal.normalize();
	} else {
		info.normal = t.gnormal;
	}

	info.dNdx = t.dNdx;
	info.dNdy = t.dNdy;

	Vector uvA = uvs[t.t[0]];
	Vector uvB = uvs[t.t[1]];
	Vector uvC = uvs[t.t[2]];

	Vector uv = uvA + (uvB - uvA) * lambda2 + (uvC - uvA) * lambda3;
	info.u = uv.x;
	info.v = uv.y;
	info.geom = this;

	return true;
}
*/
