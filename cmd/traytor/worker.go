package main

import (
	"fmt"
	"time"

	"github.com/DexterLB/traytor/rpc"
	"github.com/codegangsta/cli"
	"github.com/valyala/gorpc"
)

func runWorker(c *cli.Context) error {
	fmt.Printf(
		"will start worker with %d threads with max requests %d on this address: %s\n",
		c.GlobalInt("max-jobs"),
		c.Int("max-requests"),
		c.String("listen-address"),
	)
	address := c.String("listen-address")
	rr := rpc.NewRemoteRaytracer(
		time.Now().Unix(),
		c.GlobalInt("max-jobs"),
		c.Int("max-requests"),
		c.Int("multisample"),
	)

	w := &gorpc.Server{
		Addr:    address,
		Handler: rr.Dispatcher.NewHandlerFunc(),
	}
	if err := w.Serve(); err != nil {
		return fmt.Errorf("Cannot start rpc server: %s", err)
	}
	return nil
}
