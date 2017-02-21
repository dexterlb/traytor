traytor
=======
[![GoDoc](https://godoc.org/github.com/DexterLB/traytor?status.svg)](http://godoc.org/github.com/DexterLB/traytor)
[![Build Status](https://travis-ci.org/DexterLB/traytor.svg?branch=master)](https://travis-ci.org/DexterLB/traytor)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/DexterLB/traytor/master/LICENSE)

[![forthebadge](http://forthebadge.com/images/badges/no-ragrets.svg)](http://forthebadge.com)

```Every single ray misses```

T-ray-tor is a raytracer written in Go which uses the Path Tracing algorithm
(or something faster if we get to it)

![Skull rendered with traytor](https://github.com/DexterLB/traytor/raw/master/skull.png)

### Features

- Reads scenes from gzipped JSON (Blender export script!)
- Materials: lambert, reflective, refractive, any mixture of those
- Mesh lamps

### Usage
	$ go get github.com/DexterLB/traytor/cmd/traytor_gui

Then export your scene from Blender with the [exporter](https://github.com/DexterLB/traytor/tree/master/blender_exporter) render it with 50 samples:

	$ traytor render -t 50 my-scene.json.gz output.png

You can find some sample scenes in the sample_scenes directory.

You can also run a distributed render on many worker machines! On each worker, start:

    $ traytor worker -l :1234

to listen on port 1234, and on the client, run:

    $ traytor client -w worker1:1234 -w worker2:1234 -t 500 my-scene.json.gz output.png

this will render the scene on all workers with 500 samples.

For more info, see `traytor --help` :)
