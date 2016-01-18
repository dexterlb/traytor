package main

import (
	"log"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
)

func testSphere(display *gui.Display) {
	rnd := traytor.NewRandom(42)
	for i := 0; i < 10000; i++ {
		vec := rnd.Vec3Sphere()
		vec.Scale(250)
		vec.Add(traytor.NewVec3(400, 0, 400))
		display.SetPixel(int(vec.X), int(vec.Z), traytor.NewColour32Bit(0x0fff, 0xffff, 0x4000))
	}
}

func testHemi(display *gui.Display) {
	rnd := traytor.NewRandom(42)
	for i := 0; i < 10000; i++ {
		vec := rnd.Vec3Hemi(traytor.NewVec3(0, 1, 0))
		vec.Scale(250)
		vec.Add(traytor.NewVec3(400, 0, 400))
		display.SetPixel(int(vec.X), int(vec.Z), traytor.NewColour32Bit(0x0fff, 0xffff, 0x4000))
	}
}

func main() {
	defer gui.Quit()

	display, err := gui.NewDisplay(800, 800, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	testSphere(display)

	display.Update()

	display.Loop()
}
