package main

import (
	"github.com/DexterLB/traytor/rpc"
	"github.com/valyala/gorpc"
	"log"
)

func main() {
	rr := rpc.NewRemoteRaytracer(3)

	s := &gorpc.Server{
		Addr:    ":1234",
		Handler: rr.Dispatcher.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}
}
