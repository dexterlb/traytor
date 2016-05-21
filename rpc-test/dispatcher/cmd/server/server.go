package main

import (
	"github.com/DexterLB/traytor/rpc-test/dispatcher"
	"github.com/valyala/gorpc"
	"log"
)

type NumberList struct {
	Data []int
}

func main() {
	d := dispatcher.CreateDispatcher()
	s := &gorpc.Server{
		// Accept clients on this TCP address.
		Addr: ":1234",

		Handler: d.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}

}
