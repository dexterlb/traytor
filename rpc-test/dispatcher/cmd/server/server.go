package main

import (
	"log"

	"github.com/DexterLB/traytor/rpc-test/dispatcher"
	"github.com/valyala/gorpc"
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
	log.Printf("serving :)")
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}

}
