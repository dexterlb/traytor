package main

import (
	"image/png"
	"log"
	"os"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
)

func main() {
	display, err := gui.NewDisplay(800, 800, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	file, err := os.Open("/tmp/foo.png")
	if err != nil {
		log.Fatal(err)
		return
	}

	decoded, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	img := traytor.ToImage(decoded)

	display.ShowImage(img)

	display.Update()

	display.Loop()

	gui.Quit()
}
