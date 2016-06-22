package main

//go:generate genqrc ui

import (
	"image"
	"image/png"
	"log"
	"os"

	qml "gopkg.in/qml.v1"
)

type Interface struct {
	qml.Object
	lastImageID int
	lastImage   image.Image
}

func (i *Interface) Init(engine *qml.Engine) {
	log.Printf("QML has started!")
	engine.AddImageProvider("renderedImage", i.loadImage)
}

func (i *Interface) ShowImage(image image.Image) {
	i.lastImage = image
	i.lastImageID++
	i.Call("showImage", i.lastImageID)
}

func (i *Interface) loadImage(id string, width int, height int) image.Image {
	log.Printf("loading image")
	return i.lastImage
}

func (i *Interface) DoStuff() {
	log.Printf("doing stuff")
	pngFile, err := os.Open("/tmp/foo.png")
	if err != nil {
		log.Printf("error opening png file: %s", err)
		return
	}
	defer pngFile.Close()
	image, err := png.Decode(pngFile)
	if err != nil {
		log.Printf("error decoding png file: %s", err)
		return
	}
	i.ShowImage(image)
}

func qtLoop() error {
	engine := qml.NewEngine()

	qml.RegisterTypes("GoGui", 1, 0, []qml.TypeSpec{{
		Init: func(i *Interface, obj qml.Object) { i.Object = obj; i.Init(engine) },
	}})

	component, err := engine.LoadFile("qrc:///ui/main.qml")
	if err != nil {
		return err
	}
	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()

	return nil
}

func main() {
	err := qml.Run(qtLoop)
	if err != nil {
		log.Fatalf("can't start Qt: %s", err)
	}
}
