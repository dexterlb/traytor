package main

import (
	"log"

	"github.com/codegangsta/cli"
)

func runRender(c *cli.Context) {
	scene, image := getArguments(c)
	log.Printf(
		"will render %s to %s of size %dx%d with %d threads",
		scene, image,
		c.Int("width"), c.Int("height"),
		c.GlobalInt("max-jobs"),
	)
}
