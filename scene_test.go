package traytor

import (
	"io/ioutil"
	"testing"
)

func TestLoadSceneFromFile(t *testing.T) {
	scene, err := LoadSceneFromFile("sample_scenes/01_triangle.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Mesh.Faces) != 1 {
		t.Errorf("scene's faces should be 1, not %d", len(scene.Mesh.Faces))
	}
}

func TestLoadSceneFromBytes(t *testing.T) {
	data, err := ioutil.ReadFile("sample_scenes/01_triangle.json.gz")
	if err != nil {
		t.Fatal(err)
	}

	scene, err := LoadSceneFromBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Mesh.Faces) != 1 {
		t.Errorf("scene's faces should be 1, not %d", len(scene.Mesh.Faces))
	}
}
