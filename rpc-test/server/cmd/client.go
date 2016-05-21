package main

import (
	"github.com/DexterLB/traytor/rpc-test/server"
	"log"
	"net/rpc"
	"sync"
)

func main() {
	client, err := rpc.Dial("tcp", "192.168.0.112:1234")
	if err != nil {
		log.Fatal(err)
	}
	var reply int
	list := make([]int, 1000000)
	for i := range list {
		list[i] = 1
	}
	sumlist := &server.NumberList{Data: list}
	g := sync.WaitGroup{}
	g.Add(500)
	for i := 0; i < 500; i++ {
		go func(i int) {
			log.Printf("%d: ", i)
			err = client.Call("Sumator.Sumatorirane", sumlist, &reply)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%d", reply)
			g.Done()
		}(i)
	}
	g.Wait()
}
