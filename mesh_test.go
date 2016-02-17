package traytor

import (
	"encoding/json"
	"fmt"
)

func ExampleMesh_Intersect() {
	meshData := []byte(`{
		"vertices": [
			{
				"normal": [0, 0, 1],
				"coordinates": [0, 0, 0],
				"uv": [0, 0]
			},
			{
				"normal": [0, 0, 1],
				"coordinates": [1, 0, 0],
				"uv": [1, 0]
			},
			{
				"normal": [0, 0, 1],
				"coordinates": [1, 1, 0],
				"uv": [1, 1]
			}
		],
		"faces": [
			{
				"vertices": [0, 1, 2],
				"material": 42
			}
		]
	}`)
	mesh := &Mesh{}
	err := json.Unmarshal(meshData, &mesh)
	if err != nil {
		fmt.Printf("Error reading json: %s\n", err)
		return
	}

	ray := &Ray{
		Start:     *NewVec3(0.25, 0.75, 1),
		Direction: *NewVec3(0, 0, -1),
	}

	var (
		intersection *Intersection
		material     int
	)
	intersection, material = mesh.Intersect(ray)

	if intersection == nil {
		fmt.Printf("no intersection\n")
	} else {
		fmt.Printf("intersection point: %s\n", intersection.Point)
		fmt.Printf("caused by ray: %s -> %s\n", intersection.Incoming)
		fmt.Printf("at a distance: %.3g\n", intersection.Distance)
		fmt.Printf("with surface coordinates (%.3g, %.3g)\n",
			intersection.U, intersection.V)
		fmt.Printf("surface normal: %s\n", intersection.Normal)
		fmt.Printf("surface coordinate system: Ox: %s, Oy: %s\n",
			intersection.SurfaceOx, intersection.SurfaceOy)
		fmt.Printf("surface material: %d\n", material)
	}

	// Output:
	// intersection point: (0.25, 0.75, 0)
	// caused by ray: (0.25, 0.75, 1) -> (0, 0, -1)
	// at a distance: 1
	// with surface coordinates: (0.25, 0.75)
	// surface normal: (0, 0, 1)
	// surface coordinate system: Ox: (1, 0, 0), Oy: (0, 1, 0)
	// surface material: 42
}
