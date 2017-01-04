package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DexterLB/traytor/hdrimage"
	"github.com/codegangsta/cli"
)

func runConvert(c *cli.Context) (err error) {
	fromFilename, toFilename := getArguments(c)
	quiet := c.GlobalBool("quiet")

	if !quiet {
		log.Printf("will convert %s to %s", fromFilename, toFilename)
	}

	from, err := os.Open(fromFilename)
	if err != nil {
		return fmt.Errorf("unable to open input file: %s", err)
	}
	defer func() {
		err = from.Close()
	}()

	image, err := hdrimage.Decode(from)
	if err != nil {
		return fmt.Errorf("unable to read input image: %s", err)
	}

	fmt.Printf("first pixel: %v\n", image.Pixels[0][0])

	return saveImage(image, toFilename, "png")
}
