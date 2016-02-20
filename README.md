traytor
=======

```Every single ray misses```

T-ray-tor is a raytracer written in Go which uses the Path Tracing algorithm
(or something faster if we get to it)

### Features

- Reads scenes from JSON (Blender export script!)
- Materials: lambert, reflective, refractive, (coming soon: combinations of those)
- Mesh lamps

### Usage
	$ go get github.com/DexterLB/traytor/cmd/traytor-gui

Then export your scene from Blender with the [exporter](https://github.com/DexterLB/traytor/tree/master/blender_exporter) and run the live renderer:

	$ traytor–gui my-scene.json.gz

Note: currently the textures are loaded from the working directory, so you must be in a folder relative to the texture paths in the scene.

Soon there will be a console renderer as well.