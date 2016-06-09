package traytor

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
)

// Scene contains all the information for a scene
type Scene struct {
	Camera    *AnyCamera     `json:"camera"`
	Materials []*AnyMaterial `json:"materials"`
	Mesh      Mesh           `json:"mesh"`
	MaxDepth  int            `json:"max_depth"`
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

// LoadSceneFromFile loads the scene from a gzipped json file
func LoadSceneFromFile(filename string) (*Scene, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadScene(f)
}

// LoadSceneFromBytes loads the scene from a gzipped json byte array
func LoadSceneFromBytes(data []byte) (*Scene, error) {
	return LoadScene(bytes.NewReader(data))
}

// LoadScene loads the scene from a reader which outputs gzipped json data
func LoadScene(reader io.Reader) (*Scene, error) {
	gzReader, err := gzip.NewReader(reader)
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
	if s.MaxDepth < 1 {
		s.MaxDepth = 5
	}
}
