package scene

import (
	"io/ioutil"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	scene, err := LoadFromFile("../sample_scenes/01_triangle.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Mesh.Faces) != 1 {
		t.Errorf("scene's faces should be 1, not %d", len(scene.Mesh.Faces))
	}
}

func TestLoadFromBytes(t *testing.T) {
	data, err := ioutil.ReadFile("../sample_scenes/01_triangle.json.gz")
	if err != nil {
		t.Fatal(err)
	}

	scene, err := LoadFromBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Mesh.Faces) != 1 {
		t.Errorf("scene's faces should be 1, not %d", len(scene.Mesh.Faces))
	}
}
