package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"

	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/rpc"
)

func RenderLoop(
	sampleCounter *rpc.SampleCounter,
	client *rpc.RemoteRaytracerCaller,
	renderedImages chan<- *hdrimage.Image,
	globalSettings *rpc.SampleSettings,
	synchronous bool,
) {
	settings := *globalSettings
	for {
		settings.SamplesAtOnce = sampleCounter.Dec(globalSettings.SamplesAtOnce)
		if settings.SamplesAtOnce == 0 {
			return
		}
		var err error
		var image *hdrimage.Image

		if synchronous {
			image, err = client.Sample(&settings)
		} else {
			err = client.StoreSample(&settings)
		}
		//both samples log the error, but only synchronous sample pushes
		//image into the channel
		if err == nil {
			if synchronous {
				renderedImages <- image
			}
			log.Printf("Rendered samples :)\n")
		} else {
			log.Printf("No sample :( %s\n", err)
		}
	}
}

func JoinSamples(
	renderedImages <-chan *hdrimage.Image,
	width int,
	height int,
) *hdrimage.Image {
	averageImage := hdrimage.New(width, height)
	averageImage.Divisor = 0
	for image := range renderedImages {
		averageImage.Add(image)
		averageImage.Divisor += image.Divisor
	}
	return averageImage
}

func runClient(c *cli.Context) error {
	scene, image := getArguments(c)
	workerAdresses := c.StringSlice("worker")

	if len(workerAdresses) == 0 {
		showError(c, "can't render on zero workers :(")
	}
	synchronous := c.Bool("synchronous")

	fmt.Printf(
		"will render %s to %s of size %dx%d on those workers: %s, synchronous %v\n",
		scene, image,
		c.Int("width"), c.Int("height"),
		strings.Join(workerAdresses, ", "),
		synchronous,
	)

	width, height := c.Int("width"), c.Int("height")
	totalSamples := c.Int("total-samples")
	sampleCounter := rpc.NewSampleCounter(totalSamples)
	renderedImages := make(chan *hdrimage.Image, len(workerAdresses))
	workers := make([]*rpc.RemoteRaytracerCaller, len(workerAdresses))
	finishWorker := &sync.WaitGroup{}
	data, err := ioutil.ReadFile(scene)
	if err != nil {
		return fmt.Errorf("Error when loading scene: %s", err)
	}
	for i := range workerAdresses {
		workers[i] = rpc.NewRemoteRaytracerCaller(workerAdresses[i], 10*time.Minute)
		finishWorker.Add(1)

		requests, err := workers[i].MaxRequestsAtOnce()
		if err != nil || requests < 1 {
			return fmt.Errorf("Can't get worker's allowed requests: %s", err)
		}

		samples, err := workers[i].MaxSamplesAtOnce()
		if err != nil || samples < 1 {
			return fmt.Errorf("Can't get worker's allowed samples: %s", err)
		}

		settings := &rpc.SampleSettings{
			Width:         width,
			Height:        height,
			SamplesAtOnce: samples,
		}

		err = workers[i].LoadScene(data)
		if err != nil {
			return fmt.Errorf("Can't load scene: %s", err)
		}

		go func(i int) {
			finishRender := &sync.WaitGroup{}
			finishRender.Add(requests)
			for request := 0; request < requests; request++ {
				go func() {
					RenderLoop(sampleCounter, workers[i], renderedImages, settings, synchronous)
					finishRender.Done()
				}()
			}
			finishRender.Wait()
			if !synchronous {
				fmt.Printf("Getting image from %s\n", workerAdresses[i])
				image, err := workers[i].GetImage()
				if err == nil {
					renderedImages <- image
				} else {
					fmt.Printf("Smelly image x_x %s\n", err)
				}
			}
			finishWorker.Done()
		}(i)
	}

	go func() {
		finishWorker.Wait()
		close(renderedImages)
	}()

	averageImage := JoinSamples(renderedImages, width, height)
	file, err := os.Create(image)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("Error when saving image: %s", err)
	}
	png.Encode(file, averageImage)
	return nil
}
