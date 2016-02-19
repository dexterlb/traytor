package traytor

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetBoundingBox(t *testing.T) {
	meshData := []byte(`{
		"vertices": [
			{
				"normal": [0, 0, 1],
				"coordinates": [-1, -2, 0]
			},
			{
				"normal": [0, 0, 1],
				"coordinates": [1, 5, -8]
			},
			{
				"normal": [0, 0, 1],
				"coordinates": [1, 2, 0]
			},
			{
				"normal": [0, 0, 1],
				"coordinates": [1, 1, 0]
			}
		],
		"faces": [
			{
				"vertices": [0, 1, 2],
				"material": 42
			},
			{
				"vertices": [1, 2, 3],
				"material": 42
			}
		]
	}`)
	mesh := &Mesh{}
	err := json.Unmarshal(meshData, &mesh)
	if err != nil {
		t.Fatalf("Error reading json: %s\n", err)
		return
	}

	mesh.Init()

	bbox := mesh.GetBoundingBox()
	asserEqualVectors(t, NewVec3(1, 5, 0), &bbox.MaxVolume)
	asserEqualVectors(t, NewVec3(-1, -2, -8), &bbox.MinVolume)
}

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
				"coordinates": [0, 1, 0],
				"uv": [0, 1]
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

	mesh.Init()

	ray := &Ray{
		Start:     *NewVec3(0.15, 0.11, 1),
		Direction: *NewVec3(0, 0, -1),
	}

	var (
		intersection *Intersection
	)
	intersection = mesh.Intersect(ray)

	if intersection == nil {
		fmt.Printf("no intersection\n")
	} else {
		fmt.Printf("intersection point: %s\n", intersection.Point)
		fmt.Printf("caused by ray: %s\n", intersection.Incoming)
		fmt.Printf("at a distance: %.3g\n", intersection.Distance)
		fmt.Printf("with surface coordinates: (%.3g, %.3g)\n",
			intersection.U, intersection.V)
		fmt.Printf("surface normal: %s\n", intersection.Normal)
		fmt.Printf("surface coordinate system: Ox: %s, Oy: %s\n",
			intersection.SurfaceOx, intersection.SurfaceOy)
		fmt.Printf("surface material: %d\n", intersection.Material)
	}

	// Output:
	// intersection point: (0.15, 0.11, 0)
	// caused by ray: (0.15, 0.11, 1) -> (0, 0, -1)
	// at a distance: 1
	// with surface coordinates: (0.15, 0.11)
	// surface normal: (0, 0, 1)
	// surface coordinate system: Ox: (1, 0, 0), Oy: (0, 1, 0)
	// surface material: 42
}

func TestIntersectTwoTriangles(t *testing.T) {
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
				"coordinates": [0, 1, 0],
				"uv": [0, 1]
			},
			{
				"normal": [0, 0, -1],
				"coordinates": [0, 0, 4],
				"uv": [0, 0]
			},
			{
				"normal": [0, 0, -1],
				"coordinates": [1, 0, 4],
				"uv": [1, 0]
			},
			{
				"normal": [0, 0, -1],
				"coordinates": [0, 1, 4],
				"uv": [0, 1]
			}
		],
		"faces": [
			{
				"vertices": [3, 4, 5],
				"material": 5
			},
			{
				"vertices": [0, 1, 2],
				"material": 42
			}
		]
	}`)
	mesh := &Mesh{}
	err := json.Unmarshal(meshData, &mesh)
	if err != nil {
		t.Fatalf("Error reading json: %s\n", err)
	}

	mesh.Init()

	ray := &Ray{
		Start:     *NewVec3(0.15, 0.11, 1),
		Direction: *NewVec3(0, 0, -1),
	}

	intersection := mesh.Intersect(ray)
	if intersection == nil {
		t.Fatal("Intersection shouldn't be nil")
	}
	if intersection.Material != 42 {
		t.Error("Intersected wrong triangle")
	}

	ray.Direction.Z = 1

	intersection = mesh.Intersect(ray)
	if intersection == nil {
		t.Fatal("Intersection shouldn't be nil")
	}
	if intersection.Material != 5 {
		t.Error("Intersected wrong triangle")
	}

}

func TestIntersectWeirdTriangle(t *testing.T) {
	meshData := []byte(`{
        "faces": [
            {
                "material": 2,
                "normal": [
                    -0.008594322018325329,
                    -0.44298017024993896,
                    0.8964902758598328
                ],
                "vertices": [
                    0,
                    1,
                    2
                ]
            }
        ],
        "vertices": [
            {
                "coordinates": [
                    0.5122736096382141,
                    -0.5653880834579468,
                    0.34105977416038513
                ],
                "normal": [
                    -0.008575701154768467,
                    -0.442976176738739,
                    0.8964812159538269
                ]
            },
            {
                "coordinates": [
                    -0.9581138491630554,
                    0.14524225890636444,
                    0.6781054139137268
                ],
                "normal": [
                    -0.008575701154768467,
                    -0.442976176738739,
                    0.8964812159538269
                ]
            },
            {
                "coordinates": [
                    0.1552419662475586,
                    0.3149992525577545,
                    0.7726602554321289
                ],
                "normal": [
                    -0.008575701154768467,
                    -0.442976176738739,
                    0.8964812159538269
                ]
            }
        ]
	}`)

	cameraData := []byte(`{
        "bottom_left": [
            -0.6535904407501221,
            -1.9059882164001465,
            2.5783398151397705
        ],
        "bottom_right": [
            -0.4707377851009369,
            -1.9072681665420532,
            2.5783398151397705
        ],
        "focus": [
            -0.5628088116645813,
            -1.998731255531311,
            2.7631680965423584
        ],
        "top_left": [
            -0.6530462503433228,
            -1.8282558917999268,
            2.6456971168518066
        ],
        "top_right": [
            -0.47019362449645996,
            -1.829535961151123,
            2.6456968784332275
        ],
        "type": "pinhole"
    }`)

	mesh := &Mesh{}
	err := json.Unmarshal(meshData, &mesh)
	if err != nil {
		t.Fatalf("Error reading json: %s\n", err)
	}

	mesh.Init()

	camera := &PinholeCamera{}
	err = json.Unmarshal(cameraData, &camera)
	if err != nil {
		t.Fatalf("Error reading json: %s\n", err)
	}

	ray := camera.ShootRay(600.0/800.0, 230.0/450.0)

	intersection := mesh.Intersect(ray)
	if intersection == nil {
		t.Fatal("Intersection shouldn't be nil")
	}
	if intersection.Material != 2 {
		t.Error("Intersected wrong triangle")
	}
}
