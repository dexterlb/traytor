package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	screen *sdl.Surface
	format *sdl.PixelFormat
	window *sdl.Window
}

func (d *Display) setPixel(x int, y int, colour color.Color) {
	r, g, b, _ := colour.RGBA()

	var value uint32 = sdl.MapRGB(
		d.format,
		uint8(r>>8),
		uint8(g>>8),
		uint8(b>>8),
	)

	var pix = uintptr(unsafe.Pointer(&d.screen.Pixels()[0]))
	pix += uintptr(((y * int(d.screen.W)) + x)) * unsafe.Sizeof(value)
	var pu = unsafe.Pointer(pix)
	var pp *uint32 = (*uint32)(pu)
	*pp = value
}

func (d *Display) showImage(img image.Image) {
	topleftX := img.Bounds().Min.X
	topleftY := img.Bounds().Min.Y
	width := img.Bounds().Max.X - topleftX - 1
	height := img.Bounds().Max.Y - topleftY - 1

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			d.setPixel(x, y, img.At(x+topleftX, y+topleftY))
		}
	}
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

	display.showImage(decoded)

	display.Update()

	WaitForExit()

	sdl.Quit()
}
