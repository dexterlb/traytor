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
	gorpc.RegisterType(&traytor.Image{})
	gorpc.RegisterType(&SampleSettings{})
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

func (rr *RemoteRaytracer) StoreSample(settings *SampleSettings) {
	rr.Raytracer.StoreSample(settings)
}

func (rr *RemoteRaytracer) GetImage() *traytor.Image {
	return rr.Raytracer.GetImage()
}
