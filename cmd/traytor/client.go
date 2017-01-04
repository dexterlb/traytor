package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"

	"github.com/DexterLB/mvm/progress"
	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/rpc"
)

// RenderLoop renders samples of an image until the sample counter reaches
// 0, and then returns the combined result. If synchronous is set, images
// will be transferred from the worker after every sample. Otherwise,
// the image is transferred at the end.
func RenderLoop(
	sampleCounter *rpc.SampleCounter,
	client *rpc.RemoteRaytracerCaller,
	renderedImages chan<- *hdrimage.Image,
	globalSettings *rpc.SampleSettings,
	synchronous bool,
	bar *progress.ProgressBar,
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
			if bar != nil {
				bar.Add(1)
			}
		} else {
			log.Printf("No sample :( %s\n", err)
		}
	}
}

// JoinSamples combines samples into a single image
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

	quiet := c.GlobalBool("quiet")

	if !quiet {
		fmt.Printf(
			"will render %s to %s of size %dx%d on those workers: %s, synchronous %v\n",
			scene, image,
			c.Int("width"), c.Int("height"),
			strings.Join(workerAdresses, ", "),
			synchronous,
		)
	}

	width, height := c.Int("width"), c.Int("height")
	totalSamples := c.Int("total-samples")
	sampleCounter := rpc.NewSampleCounter(totalSamples)
	renderedImages := make(chan *hdrimage.Image, len(workerAdresses))
	workers := make([]*rpc.RemoteRaytracerCaller, len(workerAdresses))
	finishWorker := &sync.WaitGroup{}
	data, err := ioutil.ReadFile(scene)

	var bar *progress.ProgressBar
	if !quiet {
		bar = progress.StartProgressBar(totalSamples, "rendering samples ")
	}

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
					RenderLoop(sampleCounter, workers[i], renderedImages, settings, synchronous, bar)
					finishRender.Done()
				}()
			}
			finishRender.Wait()
			if !synchronous {
				if !quiet {
					bar.Prefix(fmt.Sprintf("Getting image from %s", workerAdresses[i]))
				}
				image, err := workers[i].GetImage()
				if err == nil {
					if image.Width != 0 {
						renderedImages <- image
					}
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
		if !quiet {
			bar.Done()
		}
	}()

	averageImage := JoinSamples(renderedImages, width, height)
	return saveImage(averageImage, image, c.String("format"))
}
