package traytor

import "testing"

func TestLoadScene(t *testing.T) {
	scene, err := LoadScene("sample_scenes/01_triangle.json.gz")
	if err != nil {
		t.Error(err)
	}
	t.Error(scene.Mesh.Faces)
	if len(scene.Mesh.Faces) != 1 {
		t.Errorf("scene's faces should be 1, not %d", len(scene.Mesh.Faces))
	}
}
