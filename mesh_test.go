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

	var intersection *Intersection
	intersection = mesh.Intersect(ray)
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
                "material": 0,
                "normal": [
                    0.487057089805603,
                    -0.5614719390869141,
                    -0.6689728498458862
                ],
                "vertices": [
                    0,
                    1,
                    2
                ]
            },
            {
                "material": 2,
                "normal": [
                    0.487057089805603,
                    -0.5614718198776245,
                    -0.6689728498458862
                ],
                "vertices": [
                    3,
                    4,
                    5
                ]
            }
        ],
        "vertices": [
            {
                "coordinates": [
                    0.7013661861419678,
                    0.511317253112793,
                    -1.7721754312515259
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            },
            {
                "coordinates": [
                    4.5777506828308105,
                    -0.723259449005127,
                    2.086278200149536
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            },
            {
                "coordinates": [
                    -2.290945529937744,
                    -3.9611659049987793,
                    -0.19700384140014648
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            },
            {
                "coordinates": [
                    -0.18452918529510498,
                    1.7407333850860596,
                    0.6997585296630859
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            },
            {
                "coordinates": [
                    3.6918554306030273,
                    0.5061566829681396,
                    4.5582122802734375
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            },
            {
                "coordinates": [
                    -3.1768407821655273,
                    -2.7317497730255127,
                    2.274930000305176
                ],
                "normal": [
                    0.48704490065574646,
                    -0.5614489912986755,
                    -0.6689656972885132
                ]
            }
        ]
	}`)

	cameraData := []byte(`{
        "bottom_left": [
            4.3385910987854,
            -7.2264556884765625,
            5.38138484954834
        ],
        "bottom_right": [
            4.4941935539245605,
            -7.130411148071289,
            5.38138484954834
        ],
        "focus": [
            4.4919233322143555,
            -7.300802230834961,
            5.529593467712402
        ],
        "top_left": [
            4.310408115386963,
            -7.180796146392822,
            5.469137668609619
        ],
        "top_right": [
            4.466011047363281,
            -7.084752082824707,
            5.469137668609619
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

	ray := camera.ShootRay(180.0/800.0, 355.0/450.0)

	intersection := mesh.Intersect(ray)
	if intersection == nil {
		t.Fatal("Intersection shouldn't be nil")
	}
	if intersection.Material != 0 {
		t.Error("Intersected wrong triangle")
	}
}

func TestIntersectOverlappingTriangles(t *testing.T) {
	meshData := []byte(`{
        "faces": [
            {
                "material": 2,
                "normal": [
                    0.5816099047660828,
                    -0.6875652074813843,
                    -0.4347230792045593
                ],
                "vertices": [
                    0,
                    1,
                    2
                ]
            },
            {
                "material": 2,
                "normal": [
                    0.3495648205280304,
                    0.43769627809524536,
                    -0.8283877968788147
                ],
                "vertices": [
                    3,
                    4,
                    5
                ]
            }
        ],
        "vertices": [
            {
                "coordinates": [
                    0.9402992725372314,
                    0.2813374996185303,
                    -1.3253726959228516
                ],
                "normal": [
                    0.5815912485122681,
                    -0.6875514984130859,
                    -0.4347056448459625
                ]
            },
            {
                "coordinates": [
                    3.748452663421631,
                    -0.38281798362731934,
                    3.482056140899658
                ],
                "normal": [
                    0.5815912485122681,
                    -0.6875514984130859,
                    -0.4347056448459625
                ]
            },
            {
                "coordinates": [
                    -2.6538453102111816,
                    -3.7354795932769775,
                    0.21913623809814453
                ],
                "normal": [
                    0.5815912485122681,
                    -0.6875514984130859,
                    -0.4347056448459625
                ]
            },
            {
                "coordinates": [
                    -1.213113784790039,
                    -1.5030741691589355,
                    -0.6865041255950928
                ],
                "normal": [
                    0.3495590090751648,
                    0.4376659393310547,
                    -0.8283638954162598
                ]
            },
            {
                "coordinates": [
                    1.2393922805786133,
                    2.453542709350586,
                    2.438972234725952
                ],
                "normal": [
                    0.3495590090751648,
                    0.4376659393310547,
                    -0.8283638954162598
                ]
            },
            {
                "coordinates": [
                    3.4325077533721924,
                    -4.627256870269775,
                    -0.3768632411956787
                ],
                "normal": [
                    0.3495590090751648,
                    0.4376659393310547,
                    -0.8283638954162598
                ]
            }
        ]
	}`)

	cameraData := []byte(`{
        "bottom_left": [
            4.3385910987854,
            -7.2264556884765625,
            5.38138484954834
        ],
        "bottom_right": [
            4.4941935539245605,
            -7.130411148071289,
            5.38138484954834
        ],
        "focus": [
            4.4919233322143555,
            -7.300802230834961,
            5.529593467712402
        ],
        "top_left": [
            4.310408115386963,
            -7.180796146392822,
            5.469137668609619
        ],
        "top_right": [
            4.466011047363281,
            -7.084752082824707,
            5.469137668609619
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

	ray := camera.ShootRay(520.0/800.0, 200.0/450.0)

	intersection := mesh.Intersect(ray)
	if intersection == nil {
		t.Fatal("Intersection shouldn't be nil")
	}
	if intersection.Material != 2 {
		t.Error("Intersected wrong triangle")
	}
}
