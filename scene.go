package traytor

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

// Scene contains all the information for a scene
type Scene struct {
	Camera    *AnyCamera     `json:"camera"`
	Materials []*AnyMaterial `json:"materials"`
	Mesh      Mesh           `json:"mesh"`
}

// Intersection represents a point on a surface struck by a ray
type Intersection struct {
	Point     *Vec3
	Incoming  *Ray
	Material  int
	Distance  float64
	U, V      float64
	Normal    *Vec3
	SurfaceOx *Vec3
	SurfaceOy *Vec3
}

// LoadScene loads the scene from a gzipped json file
func LoadScene(filename string) (*Scene, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	decoder := json.NewDecoder(gzReader)

	scene := &Scene{}
	err = decoder.Decode(&scene)
	if err != nil {
		return nil, err
	}
	return scene, nil
}

// Init performs all necessary preprocessing on the scene
func (s *Scene) Init() {
	s.Mesh.Init()
}
