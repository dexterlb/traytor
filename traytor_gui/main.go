package main

import (
	"image/color"
	"log"
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	screen *sdl.Surface
	format *sdl.PixelFormat
	window *sdl.Window
}

func (d *Display) setPixel(x int32, y int32, colour color.Color) {
	r, g, b, _ := colour.RGBA()

	var value uint32 = sdl.MapRGB(
		d.format,
		uint8(r),
		uint8(g),
		uint8(b),
	)

	var pix = uintptr(unsafe.Pointer(&d.screen.Pixels()[0]))
	pix += uintptr(((y * d.screen.W) + x)) * unsafe.Sizeof(value)
	var pu = unsafe.Pointer(pix)
	var pp *uint32 = (*uint32)(pu)
	*pp = value
}

func NewDisplay(width, height int, title string) (*Display, error) {
	d := &Display{}

	sdl.Init(sdl.INIT_EVERYTHING)

	var err error
	d.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width, height,
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		return nil, err
	}

	d.screen, err = d.window.GetSurface()
	if err != nil {
		return nil, err
	}

	formatID, err := d.window.GetPixelFormat()
	if err != nil {
		return nil, err
	}

	d.format, err = sdl.AllocFormat(uint(formatID))
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Display) Close() {
	d.window.Destroy()
}

func (d *Display) Update() {
	d.window.UpdateSurface()
}

type LameColour struct {
	r, g, b uint32
}

func (l *LameColour) RGBA() (uint32, uint32, uint32, uint32) {
	return l.r, l.g, l.b, 1
}

func WaitForExit() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyUpEvent:
				if t.Keysym.Sym == sdl.K_ESCAPE {
					return
				}
			}
		}
	}
}

func main() {
	display, err := NewDisplay(800, 800, "shite")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer display.Close()

	colour := &LameColour{4, 240, 20}

	for t := 0.0; t < 2*math.Pi; t += 0.001 {
		x := 16 * math.Pow(math.Sin(t), 3)
		y := 13*math.Cos(t) - 5*math.Cos(2*t) - 2*math.Cos(3*t) - math.Cos(4*t)
		display.setPixel(int32(400+x*16), int32(400-y*16), colour)
	}
	display.Update()

	WaitForExit()

	sdl.Quit()
}
