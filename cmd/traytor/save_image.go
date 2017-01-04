package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/DexterLB/traytor/hdrimage"
)

func saveImage(image *hdrimage.Image, filename string, format string) (err error) {
	file, err := os.Create(filename)
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		return fmt.Errorf("Error when saving image: %s", err)
	}

	switch format {
	case "png":
		err = png.Encode(file, image)
		if err != nil {
			return fmt.Errorf("Cannot encode png data: %s", err)
		}
	case "traytor_hdr":
		err = image.Encode(file)
		if err != nil {
			return fmt.Errorf("Cannot encode traytor_hdr data: %s", err)
		}
	default:
		return fmt.Errorf("Unknown format: %s", format)
	}

	return nil
}
