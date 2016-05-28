package rpc

import (
	"fmt"

	"github.com/DexterLB/traytor"
	"github.com/valyala/gorpc"
)

// RemoteRaytracer represents a remote raytracer and a dispatcher
type RemoteRaytracer struct {
	Scene           *traytor.Scene
	RandomGenerator *traytor.Random
	Dispatcher      *gorpc.Dispatcher
	Locker          chan struct{}
	Requests        int
	Samples         int
}

func NewRemoteRaytracer(randomSeed int64, cores int, requests int, samples int) *RemoteRaytracer {
	rr := &RemoteRaytracer{
		Scene:           nil,
		Dispatcher:      gorpc.NewDispatcher(),
		Requests:        requests,
		Locker:          make(chan struct{}, cores),
		RandomGenerator: traytor.NewRandom(randomSeed),
		Samples:         samples,
	}

	rr.Dispatcher.AddFunc("LoadScene", rr.LoadScene)
	rr.Dispatcher.AddFunc("Sample", rr.Sample)
	rr.Dispatcher.AddFunc("RequestsNumber", rr.RequestsNumber)
	gorpc.RegisterType(&traytor.Image{})
	return rr
}

func (rr *RemoteRaytracer) LoadScene(data []byte) error {
	var err error
	rr.Scene, err = traytor.LoadSceneFromBytes(data)
	rr.Scene.Init()
	return err
}

func (rr *RemoteRaytracer) Sample(size [2]int) (*traytor.Image, error) {
	rr.Locker <- struct{}{}
	defer func() { <-rr.Locker }()
	raytracer := &traytor.Raytracer{
		Scene:    rr.Scene,
		Random:   traytor.NewRandom(rr.RandomGenerator.NewSeed()),
		MaxDepth: 10,
	}
	image := traytor.NewImage(size[0], size[1])
	image.Divisor = 0
	if rr.Scene == nil {
		return nil, fmt.Errorf("Empty scene")
	}

	for i := 0; i < rr.Samples; i++ {
		raytracer.Sample(image)
	}

	return image, nil
}

func (rr *RemoteRaytracer) RequestsNumber() (int, error) {
	return rr.Requests, nil
}
