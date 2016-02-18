package main

import (
	"log"
	"os"
	"time"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()
	defer gui.Quit()

	if len(os.Args) <= 1 {
		log.Fatal("Please supply a scene file (.json.gz) as an argument.")
	}

	scene, err := traytor.LoadScene(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scene.Init()

	raytracer := &traytor.Raytracer{
		Scene:    scene,
		Random:   traytor.NewRandom(42),
		MaxDepth: 4,
	}

	width := 800
	height := 450

	display, err := gui.NewDisplay(width, height, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	image := traytor.NewImage(width, height)
	image.Divisor = 0

	for i := 0; true; i++ {
		startTime := time.Now()
		raytracer.Sample(image)
		log.Printf("rendered sample %d for %s\n", i, time.Since(startTime))
		display.ShowImage(0, 0, image)
		display.Update()
		if display.CheckExit() {
			return
		}
	}
}
