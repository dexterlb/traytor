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

	"github.com/DexterLB/traytor"
	"github.com/DexterLB/traytor/rpc"
)

type SampleCounter struct {
	sync.Mutex
	Counter int
}

func NewSampleCounter(value int) *SampleCounter {
	return &SampleCounter{Counter: value}
}

func (sc *SampleCounter) Dec(value int) int {
	sc.Lock()
	defer sc.Unlock()

	if sc.Counter >= value {
		sc.Counter -= value
	} else {
		value = sc.Counter
		sc.Counter = 0
	}
	return value
}

func RenderLoop(
	sampleCounter *SampleCounter,
	client *rpc.RemoteRaytracerCaller,
	renderedImages chan<- *traytor.Image,
	globalSettings *rpc.SampleSettings,
) {
	settings := *globalSettings

	for {
		settings.SamplesAtOnce = sampleCounter.Dec(globalSettings.SamplesAtOnce)
		if settings.SamplesAtOnce == 0 {
			return
		}

		image, err := client.Sample(&settings)
		if err == nil {
			renderedImages <- image
			log.Printf("Rendered %d samples :)\n", settings.SamplesAtOnce)
		} else {
			log.Printf("No sample :( %s\n", err)
		}
	}
}

func JoinSamples(
	renderedImages <-chan *traytor.Image,
	width int,
	height int,
) *traytor.Image {
	averageImage := traytor.NewImage(width, height)
	averageImage.Divisor = 0
	for image := range renderedImages {
		averageImage.Add(image)
		averageImage.Divisor += image.Divisor
	}
	return averageImage
}

func runClient(c *cli.Context) error {
	scene, image := getArguments(c)
	workers := c.StringSlice("worker")
	if len(workers) == 0 {
		showError(c, "can't render on zero workers :(")
	}

	fmt.Printf(
		"will render %s to %s of size %dx%d on those workers: %s\n",
		scene, image,
		c.Int("width"), c.Int("height"),
		strings.Join(workers, ", "),
	)

	width, height := c.Int("width"), c.Int("height")
	sampleCounter := NewSampleCounter(20)
	renderedImages := make(chan *traytor.Image, 20)
	clients := make([]*rpc.RemoteRaytracerCaller, len(workers))

	data, err := ioutil.ReadFile(scene)
	if err != nil {
		return fmt.Errorf("Error when loading scene: %s", err)
	}

	finishRender := &sync.WaitGroup{}

	for i := range workers {
		clients[i] = rpc.NewRemoteRaytracerCaller(workers[i], 10*time.Minute)

		requests, err := clients[i].MaxRequestsAtOnce()
		if err != nil || requests < 1 {
			return fmt.Errorf("Can't get worker's allowed requests: %s", err)
		}
		finishRender.Add(requests)

		samples, err := clients[i].MaxSamplesAtOnce()
		if err != nil || samples < 1 {
			return fmt.Errorf("Can't get worker's allowed samples: %s", err)
		}

		settings := &rpc.SampleSettings{
			Width:         width,
			Height:        height,
			SamplesAtOnce: samples,
		}

		err = clients[i].LoadScene(data)
		if err != nil {
			return fmt.Errorf("Can't load scene: %s", err)
		}

		for request := 0; request < requests; request++ {
			go func() {
				RenderLoop(sampleCounter, clients[i], renderedImages, settings)
				finishRender.Done()
			}()
		}
	}

	go func() {
		finishRender.Wait()
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
