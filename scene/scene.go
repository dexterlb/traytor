package scene

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"

	"github.com/DexterLB/traytor/camera"
	"github.com/DexterLB/traytor/materials"
	"github.com/DexterLB/traytor/mesh"
)

// Scene contains all the information for a scene
type Scene struct {
	Camera    *camera.AnyCamera        `json:"camera"`
	Materials []*materials.AnyMaterial `json:"materials"`
	Mesh      mesh.Mesh                `json:"mesh"`
	MaxDepth  int                      `json:"max_depth"`
}

// LoadFromFile loads the scene from a gzipped json file
func LoadFromFile(filename string) (scene *Scene, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()

	return Load(f)
}

// LoadFromBytes loads the scene from a gzipped json byte array
func LoadFromBytes(data []byte) (*Scene, error) {
	return Load(bytes.NewReader(data))
}

// Load loads the scene from a reader which outputs gzipped json data
func Load(reader io.Reader) (scene *Scene, err error) {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = gzReader.Close()
	}()

	decoder := json.NewDecoder(gzReader)

	scene = &Scene{}
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
