package main

import (
	"github.com/DexterLB/traytor/rpc-test/server"
	"log"
	"net"
	"net/rpc"
)

func main() {
	//rpc.HandleHTTP()
	tcpaddress, err := net.ResolveTCPAddr("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	l, e := net.ListenTCP("tcp", tcpaddress)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	s := new(server.Sumator)
	rpc.Register(s)
	rpc.Accept(l)
}
