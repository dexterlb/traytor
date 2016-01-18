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

	file, err := os.Open("foo.png")
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

	display.ShowImage(400-img.Width/2, 400-img.Height/2, img)

	display.Update()

	display.Loop()

	gui.Quit()
}
