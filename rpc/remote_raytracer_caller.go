package rpc

import (
	"time"

	"github.com/DexterLB/traytor/hdrimage"
	"github.com/valyala/gorpc"
)

// RemoteRaytracerCaller is a wrapper for calling RemoteRaytracer methods
// through RPC.
type RemoteRaytracerCaller struct {
	client     *gorpc.Client
	funcClient *gorpc.DispatcherClient
	timeout    time.Duration
}

// NewRemoteRaytracerCaller initializes the wrapper, connecting to a worker
func NewRemoteRaytracerCaller(address string, timeout time.Duration) *RemoteRaytracerCaller {
	rrc := &RemoteRaytracerCaller{
		client:  &gorpc.Client{Addr: address},
		timeout: timeout,
	}
	rrc.client.Start()
	rrc.funcClient = NewDispatcher().NewFuncClient(rrc.client)

	return rrc
}

// LoadScene sends a scene to the worker
func (rrc *RemoteRaytracerCaller) LoadScene(data []byte) error {
	_, err := rrc.funcClient.CallTimeout("LoadScene", data, rrc.timeout)
	return err
}

// MaxSamplesAtOnce gets the worker's desired samples to request at once
func (rrc *RemoteRaytracerCaller) MaxSamplesAtOnce() (int, error) {
	samples, err := rrc.funcClient.CallTimeout("MaxSamplesAtOnce", nil, rrc.timeout)
	if err != nil {
		return 0, err
	}
	return samples.(int), nil
}

// MaxRequestsAtOnce gets the worker's desired maximum simoultaneous requests
func (rrc *RemoteRaytracerCaller) MaxRequestsAtOnce() (int, error) {
	requests, err := rrc.funcClient.CallTimeout("MaxRequestsAtOnce", nil, rrc.timeout)
	if err != nil {
		return 0, err
	}
	return requests.(int), nil
}

// Sample waits for the worker to sample an image, retreives it and returns it
func (rrc *RemoteRaytracerCaller) Sample(settings *SampleSettings) (*hdrimage.Image, error) {
	image, err := rrc.funcClient.CallTimeout("Sample", settings, rrc.timeout)
	if err != nil {
		return nil, err
	}
	return image.(*hdrimage.Image), nil
}

// StoreSample waits for the worker to sample an image, storing it worker-side
func (rrc *RemoteRaytracerCaller) StoreSample(settings *SampleSettings) error {
	_, err := rrc.funcClient.CallTimeout("StoreSample", settings, rrc.timeout)
	return err
}

// GetImage retreives the combined result of any previously stored samples
func (rrc *RemoteRaytracerCaller) GetImage() (*hdrimage.Image, error) {
	image, err := rrc.funcClient.CallTimeout("GetImage", nil, rrc.timeout)
	if err != nil {
		return nil, err
	}
	return image.(*hdrimage.Image), nil
}
