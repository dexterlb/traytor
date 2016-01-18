package gui

import (
	"image"
	"image/color"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// Display is a window with a drawable canvas
type Display struct {
	screen *sdl.Surface
	format *sdl.PixelFormat
	window *sdl.Window
	width  int
	height int
}

// SetPixel sets the colour of the pixel at the given coordinates
func (d *Display) SetPixel(x int, y int, colour color.Color) {
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

// ShowImage draws an image at the given coordinates
func (d *Display) ShowImage(x int, y int, img image.Image) {
	topleftX := img.Bounds().Min.X
	topleftY := img.Bounds().Min.Y
	width := img.Bounds().Max.X - topleftX - 1
	height := img.Bounds().Max.Y - topleftY - 1

	for locX := 0; locX < width; locX++ {
		for locY := 0; locY < height; locY++ {
			d.SetPixel(x+locX, y+locY, img.At(locX+topleftX, locY+topleftY))
		}
	}
}

// Fill fills the screen with the specified colour
func (d *Display) Fill(colour color.Color) {
	for x := 0; x < d.width; x++ {
		for y := 0; y < d.height; y++ {
			d.SetPixel(x, y, colour)
		}
	}
}

// NewDisplay creates a display with the given size and window title
func NewDisplay(width, height int, title string) (*Display, error) {
	d := &Display{}

	sdl.Init(sdl.INIT_EVERYTHING)

	d.width = width
	d.height = height

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

// Close closes the window of the display
func (d *Display) Close() {
	d.window.Destroy()
}

// Update refreshes the screen
func (d *Display) Update() {
	d.window.UpdateSurface()
}

// CheckExit returns true if an exit event has happened, false otherwise
func (d *Display) CheckExit() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyUpEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				return true
			}
		case *sdl.WindowEvent:
			d.Update()
		}
	}
	return false
}

// Loop waits for an exit event, refreshing the screen each time it's uncovered
func (d *Display) Loop() {
	for !d.CheckExit() {
	}
}

// Quit exits the gui
func Quit() {
	sdl.Quit()
}
