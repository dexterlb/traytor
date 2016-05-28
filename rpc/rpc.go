package rpc

import (
	"fmt"

	"github.com/DexterLB/traytor"
	"github.com/valyala/gorpc"
)

// RemoteRaytracer represents a remote raytracer and a dispatcher
type RemoteRaytracer struct {
	Raytracer  *traytor.Raytracer
	Dispatcher *gorpc.Dispatcher
	Cores      int
}

func NewRemoteRaytracer(randomSeed int64, cores int) *RemoteRaytracer {
	raytracer := &traytor.Raytracer{
		Scene:    nil,
		Random:   traytor.NewRandom(int64(randomSeed)),
		MaxDepth: 10,
	}
	dispatcher := gorpc.NewDispatcher()

	rr := &RemoteRaytracer{
		Raytracer:  raytracer,
		Dispatcher: dispatcher,
		Cores:      cores,
	}

	rr.Dispatcher.AddFunc("LoadScene", rr.LoadScene)
	rr.Dispatcher.AddFunc("Sample", rr.Sample)
	rr.Dispatcher.AddFunc("CoresNumber", rr.CoresNumber)
	gorpc.RegisterType(&traytor.Image{})
	return rr
}

func (rr *RemoteRaytracer) LoadScene(data []byte) error {
	var err error
	rr.Raytracer.Scene, err = traytor.LoadSceneFromBytes(data)
	rr.Raytracer.Scene.Init()
	return err
}

func (rr *RemoteRaytracer) Sample(size [2]int) (*traytor.Image, error) {
	image := traytor.NewImage(size[0], size[1])
	image.Divisor = 0
	if rr.Raytracer.Scene == nil {
		return nil, fmt.Errorf("errorerror error roroor")
	}
	rr.Raytracer.Sample(image)
	return image, nil
}

func (rr *RemoteRaytracer) CoresNumber() (int, error) {
	return rr.Cores
}
