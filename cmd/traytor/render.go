package main

import (
	"log"

	"github.com/codegangsta/cli"
)

func runRender(c *cli.Context) error {
	scenePath, image := getArguments(c)
	log.Printf(
		"will render %d samples of %s to %s of size %dx%d with %d threads",
		c.Int("total-samples"),
		scenePath, image,
		c.Int("width"), c.Int("height"),
		c.Int("max-jobs"),
	)

	/*
		width, height := c.Int("width"), c.Int("height")
		totalSamples := c.Int("total-samples")
		sampleCounter := rpc.NewSampleCounter(totalSamples)
		renderedImages := make(chan *hdrimage.Image)

		randomGen := random.New(42)

		scene, err := scene.LoadFromFile(scenePath)
		if err != nil {
			return fmt.Errorf("can't open scene: %s", err)
		}

		for i := 0; i < c.Int("max-jobs"); i++ {
			go func() {
				raytracer := raytracer.Raytracer{
					Scene:  scene,
					Random: random.New(randomGen.NewSeed()),
				}

				image := hdrimage.New(width, height)
				image.Divisor = 0

				for {
					// sample image and decrease counter
				}

			}()
		}
	*/
	return nil
}
