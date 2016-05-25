package main

import (
	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/rpc"
	"github.com/valyala/gorpc"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	c := &gorpc.Client{
		Addr: ":1234",
	}
	c.Start()
	rr := rpc.NewRemoteRaytracer(3)
	d := rr.Dispatcher.NewFuncClient(c)

	data, _ := ioutil.ReadFile(os.Args[1])
	_, err := d.Call("LoadScene", data)
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
	var image *traytor.Image
	image, err = d.Call("Sample", [2]int{800, 600})
	file, _ := os.Create("name")

	err = png.Encode(file, image)
	file.Close()
	if err != nil {
		log.Fatalf("Error when sending request to server: %s", err)
	}
}
