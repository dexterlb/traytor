package main

import (
	"log"
	"time"

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/gui"
)

func testSphere(display *gui.Display) {
	rnd := traytor.NewRandom(42)
	for !display.CheckExit() {
		vec := rnd.Vec3Sphere()
		vec.Scale(350)
		vec.Add(traytor.NewVec3(400, 0, 400))
		display.SetPixel(int(vec.X), int(vec.Z), traytor.NewColour32Bit(0x0fff, 0xffff, 0x4000))
		time.Sleep(50 * time.Microsecond)
		display.Update()
	}
}

func testHemi(display *gui.Display) {
	rnd := traytor.NewRandom(42)
	for !display.CheckExit() {
		vec := rnd.Vec3Hemi(traytor.NewVec3(0, 0, -1))
		vec.Scale(350)
		vec.Add(traytor.NewVec3(400, 0, 400))
		display.SetPixel(int(vec.X), int(vec.Z), traytor.NewColour32Bit(0x0fff, 0xffff, 0x4000))

		vecCos := rnd.Vec3HemiCos(traytor.NewVec3(0, 0, 1))
		vecCos.Scale(350)
		vecCos.Add(traytor.NewVec3(400, 0, 400))
		display.SetPixel(int(vecCos.X), int(vecCos.Z), traytor.NewColour32Bit(0xffff, 0x0fff, 0x4000))

		time.Sleep(50 * time.Microsecond)
		display.Update()
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

	testHemi(display)
}
