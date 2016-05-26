package main

import (
	"fmt"
	"github.com/DexterLB/traytor/rpc"
	"github.com/codegangsta/cli"
	"github.com/valyala/gorpc"
)

func runServer(c *cli.Context) error {
	fmt.Printf(
		"will start server with %d threads on this address: %s\n",
		c.GlobalInt("max-jobs"),
		c.String("listen-address"),
	)
	address := c.String("listen-address")
	rr := rpc.NewRemoteRaytracer(3)

	s := &gorpc.Server{
		Addr:    address,
		Handler: rr.Dispatcher.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		return fmt.Errorf("Cannot start rpc server: %s", err)
	}
	return nil
}
