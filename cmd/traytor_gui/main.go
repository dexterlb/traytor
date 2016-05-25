package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
	"github.com/pkg/profile"
)

func renderBucket(scene *traytor.Scene, width int, height int,
	samples chan<- *traytor.Image, threadId int) {

	raytracer := &traytor.Raytracer{
		Scene:    scene,
		Random:   traytor.NewRandom(int64(threadId)),
		MaxDepth: 10,
	}

	for {
		image := traytor.NewImage(width, height)
		startTime := time.Now()
		raytracer.Sample(image)
		log.Printf("thread %d rendered sample for %s\n",
			threadId, time.Since(startTime))

		samples <- image
	}
}

func main() {
	defer profile.Start().Stop()
	defer gui.Quit()

	if len(os.Args) <= 1 {
		log.Fatal("Please supply a scene file (.json.gz) as an argument.")
	}

	scene, err := traytor.LoadSceneFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scene.Init()

	width := 800
	height := 450

	display, err := gui.NewDisplay(width, height, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	samples := make(chan *traytor.Image, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go renderBucket(scene, width, height, samples, i)
	}

	image := traytor.NewImage(width, height)
	image.Divisor = 0

	for i := 0; true; i++ {
		sample := <-samples
		image.Add(sample)
		image.Divisor++

		log.Printf("displayed sample %d\n", i)

		display.ShowImage(0, 0, image)
		display.Update()
		if display.CheckExit() {
			return
		}
	}
}
