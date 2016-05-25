package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/rpc"
	"github.com/valyala/gorpc"
)

func main() {
	c := &gorpc.Client{
		Addr: "hoth:1234",
	}
	c.Start()
	rr := rpc.NewRemoteRaytracer(3)
	d := rr.Dispatcher.NewFuncClient(c)

	data, _ := ioutil.ReadFile(os.Args[1])
	_, err := d.Call("LoadScene", data)
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
	resp, err := d.Call("Sample", [2]int{800, 450})
	file, _ := os.Create("name")
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
	png.Encode(file, resp.(*traytor.Image))
	file.Close()
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
}
