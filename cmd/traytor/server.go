package main

import (
	"log"

	"github.com/codegangsta/cli"
)

func runServer(c *cli.Context) {
	log.Printf(
		"will start server with %d threads on this address: %s",
		c.GlobalInt("max-jobs"),
		c.String("listen-address"),
	)
}
