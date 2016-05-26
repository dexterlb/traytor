package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/rpc"
	"github.com/valyala/gorpc"
)

func runClient(c *cli.Context) error {
	scene, image := getArguments(c)
	workers := c.StringSlice("worker")
	if len(workers) == 0 {
		showError(c, "can't render on zero workers :(")
	}

	fmt.Printf(
		"will render %s to %s of size %dx%d with %d threads on those workers: %s\n",
		scene, image,
		c.Int("width"), c.Int("height"),
		c.GlobalInt("max-jobs"),
		strings.Join(workers, ", "),
	)
	width, height := c.Int("width"), c.Int("height")
	address := workers[0]
	client := &gorpc.Client{
		Addr: address,
	}
	client.Start()
	rr := rpc.NewRemoteRaytracer(3)
	d := rr.Dispatcher.NewFuncClient(client)

	data, _ := ioutil.ReadFile(scene)
	_, err := d.CallTimeout("LoadScene", data, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("Error when sending request to server: %s", err)
	}
	resp, err := d.CallTimeout("Sample", [2]int{width, height}, 10*time.Minute)
	file, _ := os.Create(image)
	if err != nil {
		return fmt.Errorf("Error when sending request to server: %s", err)
	}
	png.Encode(file, resp.(*traytor.Image))
	file.Close()
	if err != nil {
		return fmt.Errorf("Error when sending request to server: %s", err)
	}
	return nil
}
