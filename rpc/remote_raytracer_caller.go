package rpc

import (
	"time"

	"github.com/DexterLB/traytor"
	"github.com/valyala/gorpc"
)

type RemoteRaytracerCaller struct {
	client     *gorpc.Client
	funcClient *gorpc.DispatcherClient
	timeout    time.Duration
}

func NewRemoteRaytracerCaller(address string, timeout time.Duration) *RemoteRaytracerCaller {
	rrc := &RemoteRaytracerCaller{
		client:  &gorpc.Client{Addr: address},
		timeout: timeout,
	}
	rrc.client.Start()
	rrc.funcClient = NewDispatcher().NewFuncClient(rrc.client)

	return rrc
}

func (rrc *RemoteRaytracerCaller) LoadScene(data []byte) error {
	_, err := rrc.funcClient.CallTimeout("LoadScene", data, rrc.timeout)
	return err
}

func (rrc *RemoteRaytracerCaller) MaxSamplesAtOnce() (int, error) {
	samples, err := rrc.funcClient.CallTimeout("MaxSamplesAtOnce", nil, rrc.timeout)
	if err != nil {
		return 0, err
	}
	return samples.(int), nil
}

func (rrc *RemoteRaytracerCaller) MaxRequestsAtOnce() (int, error) {
	requests, err := rrc.funcClient.CallTimeout("MaxRequestsAtOnce", nil, rrc.timeout)
	if err != nil {
		return 0, err
	}
	return requests.(int), nil
}

func (rrc *RemoteRaytracerCaller) Sample(settings *SampleSettings) (*traytor.Image, error) {
	image, err := rrc.funcClient.CallTimeout("Sample", settings, rrc.timeout)
	if err != nil {
		return nil, err
	}
	return image.(*traytor.Image), nil
}
