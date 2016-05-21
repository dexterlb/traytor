package main

import (
	//"github.com/DexterLB/traytor/rpc-test/server"
	"github.com/valyala/gorpc"
	"log"
)

type NumberList struct {
	Data []int
}

func Sum(input *NumberList) int {
	reply := 0
	for i := range input.Data {
		reply = reply + input.Data[i]
	}
	return reply
}

func main() {
	d := gorpc.NewDispatcher()
	d.AddFunc("Sum", Sum)
	s := &gorpc.Server{
		// Accept clients on this TCP address.
		Addr: ":1234",

		Handler: d.NewHandlerFunc(),
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}

}
