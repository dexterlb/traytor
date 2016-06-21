package rpc

import (
	"github.com/DexterLB/traytor"
	"github.com/valyala/gorpc"
)

// RemoteRaytracer represents a remote raytracer and a dispatcher
type RemoteRaytracer struct {
	Raytracer  *ConcurrentRaytracer
	Requests   int
	Dispatcher *gorpc.Dispatcher
	Samples    int
}

// SampleSettings contains parameters for making a sample
type SampleSettings struct {
	Width         int
	Height        int
	SamplesAtOnce int
}

func NewRemoteRaytracer(
	randomSeed int64,
	threads int, // number of threads that render simoultaneously
	maxRequestsAtOnce int, // the requests we accept at once (2*threads is a good number)
	samplesAtOnce int, // number of samples to send at once to the client
) *RemoteRaytracer {
	rr := &RemoteRaytracer{
		Samples:    samplesAtOnce,
		Raytracer:  NewConcurrentRaytracer(threads, nil, randomSeed),
		Dispatcher: gorpc.NewDispatcher(),
		Requests:   maxRequestsAtOnce,
	}

	rr.Dispatcher.AddFunc("LoadScene", rr.LoadScene)
	rr.Dispatcher.AddFunc("Sample", rr.Sample)
	rr.Dispatcher.AddFunc("MaxRequestsAtOnce", rr.MaxRequestsAtOnce)
	rr.Dispatcher.AddFunc("MaxSamplesAtOnce", rr.MaxSamplesAtOnce)
	gorpc.RegisterType(&traytor.Image{})
	gorpc.RegisterType(&SampleSettings{})
	return rr
}

func (rr *RemoteRaytracer) LoadScene(data []byte) error {
	var err error
	scene, err := traytor.LoadSceneFromBytes(data)
	if err != nil {
		return err
	}
	scene.Init()
	rr.Raytracer.SetScene(scene)
	return nil
}

func (rr *RemoteRaytracer) Sample(settings *SampleSettings) (*traytor.Image, error) {
	return rr.Raytracer.Sample(settings)
}

func (rr *RemoteRaytracer) MaxRequestsAtOnce() (int, error) {
	return rr.Requests, nil
}

func (rr *RemoteRaytracer) MaxSamplesAtOnce() (int, error) {
	return rr.Samples, nil
}
