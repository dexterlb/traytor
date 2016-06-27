package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func savePng(image image.Image, filename string) (err error) {
	file, err := os.Create(filename)
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		return fmt.Errorf("Error when saving image: %s", err)
	}
	err = png.Encode(file, image)
	if err != nil {
		return fmt.Errorf("Cannot encode png data: %s", err)
	}

	return nil
}
