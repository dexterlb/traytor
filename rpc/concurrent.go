package rpc

import (
	"fmt"

	"github.com/DexterLB/traytor"
)

// ConcurrentRaytracer can render image samples on a scene in parallel,
// and has locks to ensure a maximum number of parallel renders.
// It can store samples internally, and they can be collected on demand.
type ConcurrentRaytracer struct {
	parallelSamples int
	units           chan *renderUnit
}

type renderUnit struct {
	raytracer traytor.Raytracer
	image     *traytor.Image
}

// NewConcurrentRaytracer creates a concurrent raytracer with parallelSamples
// allowed number of parallel operations
func NewConcurrentRaytracer(
	parallelSamples int,
	scene *traytor.Scene,
	seed int64,
) *ConcurrentRaytracer {
	if parallelSamples < 1 {
		panic("must have at least one rendering unit")
	}

	cr := &ConcurrentRaytracer{
		parallelSamples: parallelSamples,
		units:           make(chan *renderUnit, parallelSamples),
	}

	random := traytor.NewRandom(seed)

	for i := 0; i < parallelSamples; i++ {
		cr.units <- &renderUnit{
			raytracer: traytor.Raytracer{
				Scene:  scene,
				Random: traytor.NewRandom(random.NewSeed()),
			},
			image: nil,
		}
	}

	return cr
}

// StoreSample renders the scene into an internal image. You can (and should)
// call it multiple times, in parallel, and when you're finished you can get
// the merged samples with GetImage(). StoreSample will block if the parallel
// calls exceed the parallelSamples value, and wait for other samples to finish.
func (cr *ConcurrentRaytracer) StoreSample(settings *SampleSettings) error {
	unit := <-cr.units
	if unit.raytracer.Scene == nil {
		return fmt.Errorf("N/A scene")
	}

	if unit.image == nil {
		unit.image = traytor.NewImage(settings.Width, settings.Height)
		unit.image.Divisor = 0
	}

	for i := 0; i < settings.SamplesAtOnce; i++ {
		unit.raytracer.Sample(unit.image)
	}

	cr.units <- unit
	return nil
}

// Sample works like StoreSample(), but instead of storing the image internally,
// returns a new image.
func (cr *ConcurrentRaytracer) Sample(settings *SampleSettings) (*traytor.Image, error) {
	unit := <-cr.units

	if unit.raytracer.Scene == nil {
		return nil, fmt.Errorf("N/A scene")
	}
	image := traytor.NewImage(settings.Width, settings.Height)
	image.Divisor = 0

	for i := 0; i < settings.SamplesAtOnce; i++ {
		unit.raytracer.Sample(image)
	}

	cr.units <- unit

	return image, nil
}

func (cr *ConcurrentRaytracer) getAllUnits() []*renderUnit {
	units := make([]*renderUnit, cr.parallelSamples)
	for i := range units {
		unit := <-cr.units
		units[i] = unit
	}
	return units
}

func (cr *ConcurrentRaytracer) pushAllUnits(units []*renderUnit) {
	for i := range units {
		cr.units <- units[i]
	}
}

func (cr *ConcurrentRaytracer) SetScene(scene *traytor.Scene) {
	units := cr.getAllUnits()
	for _, unit := range units {
		unit.image = nil
		unit.raytracer.Scene = scene
	}
	fmt.Printf("%d - %d\n", len(units), cr.ParallelSamples())
	cr.pushAllUnits(units)
}

// GetImage collects all samples up to this moment (and waits for those that
// are currently being rendered to finish), and merges them.
// StoreSample() can be called during calling this function, and will block until
// it finishes. Next samples will start from zero (e.g. the base image is reset)
func (cr *ConcurrentRaytracer) GetImage() *traytor.Image {
	var mergedSamples *traytor.Image

	units := cr.getAllUnits()
	for _, unit := range units {
		if mergedSamples == nil {
			mergedSamples = unit.image
		} else if unit.image != nil {
			mergedSamples.Add(unit.image)
			mergedSamples.Divisor += unit.image.Divisor
		}
		unit.image = nil
	}
	cr.pushAllUnits(units)

	return mergedSamples
}

func (cr *ConcurrentRaytracer) ParallelSamples() int {
	return cr.parallelSamples
}
