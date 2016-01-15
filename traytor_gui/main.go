package main

import (
	"fmt"
	"image/color"

	"github.com/salviati/go-qt5/qt5"
)

func main() {
	qt5.Main(ui_main)
}

func ui_main() {
	w := qt5.NewWidget()
	defer w.Close()

	// canvas := qt5.NewWidget()

	w.OnPaintEvent(func(e *qt5.PaintEvent) {
		paint := qt5.NewPainter()
		defer paint.Close()
		paint.Begin(w)
		pen := qt5.NewPen()
		pen.SetColor(color.RGBA{255, 128, 0, 0})
		pen.SetWidth(2)
		fmt.Println(pen, pen.Color(), pen.Width())
		paint.SetPen(pen)
		brush := qt5.NewBrush()
		brush.SetStyle(qt5.SolidPattern)
		brush.SetColor(color.RGBA{128, 128, 0, 255})
		paint.SetBrush(brush)
		paint.DrawRect(qt5.Rect{10, 10, 100, 100})
	})

	w.SetSize(qt5.Size{400, 400})

	w.Show()

	qt5.Run()
}
