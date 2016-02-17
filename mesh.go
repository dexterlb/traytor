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

func (m *Mesh) IntersectionTriangle(ray *Ray, triangle *Triangle, maxDistance float64) *Intersection {
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

	if intersectDist < 0 || intersectDist > maxDistance {
		return nil
	}
	//lambda2 = (dist^dir)*AC / det
	//lambda3 = -(dist^dir)*AB / det 
	float64 lambda2 = MixedProduct(intersectDist, rayDir, minusVectors(C, A)) * reverseDet 
	float64 lambda3 = MixedProduct(intersectDist, rayDir, minusVectors(B, A)) * reverseDet
	if lambda2 < 0 || lambda2 > 1 || lambda3 < 0 || lambda3 > 1 || lambda2 + lambda3 > 1 {
		return nil
	} 
	intersection.Distance = intersectDist
	intersection.Point = ray.Start + rayDir * intersectDist
	if Triangle.Normal {
			intersection.Normal = Triangle.Normal
		} else {
			Anormal := m.Vertices[triangle.Vertices[0]].Normal
			Bnormal := m.Vertices[triangle.Vertices[1]].Normal
			Cnormal := m.Vertices[triangle.Vertices[2]].Normal
			ABxlambda2 := MinusVectors(Bnormal, Anormal).Scaled(lambda2)
			ACxlambda3 := MinusVectors(Cnormal, Anormal).Scaled(lambda3)
			intersection.Normal = AddVectors(Anormal, AddVectors(ABxlambda2, ACxlambda3))
		}
	return intersection
}

/*

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
