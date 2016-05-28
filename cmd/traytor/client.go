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
	"github.com/valyala/gorpc"
)

type SampleCounter struct {
	sync.Mutex
	Counter int
}

func NewSampleCounter(value int) *SampleCounter {
	return &SampleCounter{Counter: value}
}

func (sc *SampleCounter) Dec() bool {
	sc.Lock()
	defer sc.Unlock()

	if sc.Counter > 0 {
		sc.Counter--
		return true
	}
	return false
}

func RenderLoop(
	sampleCounter *SampleCounter,
	client *gorpc.DispatcherClient,
	renderedImages chan<- *traytor.Image,
	width int,
	height int,
) {
	for {
		if !sampleCounter.Dec() {
			return
		}
		image, err := client.CallTimeout("Sample", [2]int{width, height}, 10*time.Minute)
		if err == nil {
			renderedImages <- image.(*traytor.Image)
			log.Printf("Rendered sample :)\n")
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
		"will render %s to %s of size %dx%d with %d threads on those workers: %s\n",
		scene, image,
		c.Int("width"), c.Int("height"),
		c.GlobalInt("max-jobs"),
		strings.Join(workers, ", "),
	)
	rr := rpc.NewRemoteRaytracer(
		time.Now().Unix(),
		c.GlobalInt("max-jobs"),
		c.GlobalInt("max-jobs")*2,
		1,
	)

	width, height := c.Int("width"), c.Int("height")
	sampleCounter := NewSampleCounter(400)
	renderedImages := make(chan *traytor.Image, 400)
	clients := make([]*gorpc.Client, len(workers))

	data, err := ioutil.ReadFile(scene)
	if err != nil {
		return fmt.Errorf("Error when loading scene: %s", err)
	}

	finishRender := &sync.WaitGroup{}

	for i := range workers {
		clients[i] = &gorpc.Client{Addr: workers[i]}
		clients[i].Start()
		dispatcher := rr.Dispatcher.NewFuncClient(clients[i])
		cores, err := dispatcher.CallTimeout("RequestsNumber", nil, 10*time.Minute)
		if err != nil || cores.(int) < 1 {
			return fmt.Errorf("Problem with calculating client cores: %s", err)
		}
		finishRender.Add(cores.(int))
		_, err = dispatcher.CallTimeout("LoadScene", data, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Can't load scene: %s", err)
		}
		for core := 0; core < cores.(int); core++ {
			go func() {
				RenderLoop(sampleCounter, dispatcher, renderedImages, width, height)
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
