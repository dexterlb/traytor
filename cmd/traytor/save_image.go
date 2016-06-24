package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func savePng(image image.Image, filename string) error {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("Error when saving image: %s", err)
	}
	png.Encode(file, image)
	return nil
}
