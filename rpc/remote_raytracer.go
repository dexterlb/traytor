package rpc

import (
	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/scene"
	"github.com/valyala/gorpc"
)

// RemoteRaytracer represents a remote raytracer and a dispatcher.
// It can be used in a client/server environment with gorpc.
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

// NewRemoteRaytracer initialises the remote raytracer object
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

	rr.registerFunctions()

	return rr
}

// NewDispatcher returns only the dispatcher
// useful on the client, where you don't need the entire RemoteRaytracer
func NewDispatcher() *gorpc.Dispatcher {
	rr := &RemoteRaytracer{
		Dispatcher: gorpc.NewDispatcher(),
	}

	rr.registerFunctions()

	return rr.Dispatcher
}

func (rr *RemoteRaytracer) registerFunctions() {
	rr.Dispatcher.AddFunc("LoadScene", rr.LoadScene)
	rr.Dispatcher.AddFunc("Sample", rr.Sample)
	rr.Dispatcher.AddFunc("MaxRequestsAtOnce", rr.MaxRequestsAtOnce)
	rr.Dispatcher.AddFunc("MaxSamplesAtOnce", rr.MaxSamplesAtOnce)
	rr.Dispatcher.AddFunc("StoreSample", rr.StoreSample)
	rr.Dispatcher.AddFunc("GetImage", rr.GetImage)
	gorpc.RegisterType(&hdrimage.Image{})
	gorpc.RegisterType(&SampleSettings{})
}

// LoadScene loads a scene
func (rr *RemoteRaytracer) LoadScene(data []byte) error {
	var err error
	scene, err := scene.LoadFromBytes(data)
	if err != nil {
		return err
	}
	scene.Init()
	rr.Raytracer.SetScene(scene)
	return nil
}

// Sample samples an image and returns it
func (rr *RemoteRaytracer) Sample(settings *SampleSettings) (*hdrimage.Image, error) {
	return rr.Raytracer.Sample(settings)
}

// MaxRequestsAtOnce returns the maximum number of requests allowed to the worker
// at the same time
func (rr *RemoteRaytracer) MaxRequestsAtOnce() (int, error) {
	return rr.Requests, nil
}

// MaxSamplesAtOnce returns the number of samples the worker might render at once
func (rr *RemoteRaytracer) MaxSamplesAtOnce() (int, error) {
	return rr.Samples, nil
}

// StoreSample stores samples an image without returning it, to be used
// with a later call of GetImage()
func (rr *RemoteRaytracer) StoreSample(settings *SampleSettings) error {
	return rr.Raytracer.StoreSample(settings)
}

// GetImage returns the combined result of any previously stored samples
func (rr *RemoteRaytracer) GetImage() *hdrimage.Image {
	return rr.Raytracer.GetImage()
}
