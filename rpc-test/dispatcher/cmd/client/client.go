package main

import (
	"github.com/DexterLB/traytor/rpc-test/dispatcher"
	"github.com/valyala/gorpc"
	"log"
)

func main() {
	c := &gorpc.Client{
		// TCP address of the server.
		Addr: "192.168.0.112:1234",
	}
	c.Start()
	d := dispatcher.CreateDispatcher().NewFuncClient(c)

	resp, err := d.Call("Sum", &dispatcher.NumberList{Data: []int{1, 2, 3, 4, 5, 6}})
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
	log.Printf("%d", resp)
}
