package main

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"
)

func runClient(c *cli.Context) {
	scene, image := getArguments(c)
	servers := c.StringSlice("server")
	if len(servers) == 0 {
		showError(c, "can't render on zero servers :(")
	}

	log.Printf(
		"will render %s to %s of size %dx%d with %d threads on those servers: %v",
		scene, image,
		c.Int("width"), c.Int("height"),
		c.GlobalInt("max-jobs"),
		strings.Join(servers, ", "),
	)
}
