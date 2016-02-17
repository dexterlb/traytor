package main

import (
	"log"
	"os"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
)

func main() {
	defer gui.Quit()
	display, err := gui.NewDisplay(800, 450, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	if len(os.Args) <= 1 {
		log.Fatal("Please supply a scene file (.json.gz) as an argument.")
	}

	scene, err := traytor.LoadScene(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	scene.Init()

	raytracer := &traytor.Raytracer{
		Scene:  scene,
		Random: traytor.NewRandom(42),
	}

	image := traytor.NewImage(800, 450)
	image.Divisor = 0

	for {
		raytracer.Sample(image)
		display.ShowImage(0, 0, image)
		display.Update()
		if display.CheckExit() {
			return
		}
	}
}
