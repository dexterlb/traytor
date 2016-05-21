package main

import (
	"github.com/DexterLB/traytor/rpc-test/server"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	var reply int
	list := make([]int, 100000000)
	for i := range list {
		list[i] = 1
	}
	log.Printf("%s", "reply")
	sumlist := &server.NumberList{Data: list}
	err = client.Call("Sumator.Sumatorirane", sumlist, &reply)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d", reply)
}
