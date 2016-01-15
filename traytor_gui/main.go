package main

import (
	"io/ioutil"
	"log"

	"github.com/salviati/go-qt5/qt5"
)

func main() {
	qt5.Main(ui_main)
}

func ui_main() {
	data, err := ioutil.ReadFile("/tmp/foo.png")
	if err != nil {
		log.Fatal(err)
		return
	}

	pixmap := qt5.NewPixmapWithData(data)

	w := qt5.NewWidget()
	defer w.Close()

	w.SetWindowTitle("This is a big shit")
	w.SetSizev(640, 480)

	label := qt5.NewLabel()
	label.SetPixmap(pixmap)

	vbox := qt5.NewVBoxLayout()
	vbox.AddWidget(label)

	w.SetLayout(vbox)
	w.Show()

	qt5.Run()
}
