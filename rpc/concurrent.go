package rpc

import (
	"fmt"

	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/random"
	"github.com/DexterLB/traytor/raytracer"
	"github.com/DexterLB/traytor/scene"
)

// ConcurrentRaytracer can render image samples on a scene in parallel,
// and has locks to ensure a maximum number of parallel renders.
// It can store samples internally, and they can be collected on demand.
type ConcurrentRaytracer struct {
	parallelSamples int
	units           chan *renderUnit
}

type renderUnit struct {
	raytracer raytracer.Raytracer
	image     *hdrimage.Image
}

// NewConcurrentRaytracer creates a concurrent raytracer with parallelSamples
// allowed number of parallel operations
func NewConcurrentRaytracer(
	parallelSamples int,
	scene *scene.Scene,
	seed int64,
) *ConcurrentRaytracer {
	if parallelSamples < 1 {
		panic("must have at least one rendering unit")
	}

	cr := &ConcurrentRaytracer{
		parallelSamples: parallelSamples,
		units:           make(chan *renderUnit, parallelSamples),
	}

	randomGen := random.New(seed)

	for i := 0; i < parallelSamples; i++ {
		cr.units <- &renderUnit{
			raytracer: raytracer.Raytracer{
				Scene:  scene,
				Random: random.New(randomGen.NewSeed()),
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
		unit.image = hdrimage.New(settings.Width, settings.Height)
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
func (cr *ConcurrentRaytracer) Sample(settings *SampleSettings) (*hdrimage.Image, error) {
	unit := <-cr.units

	if unit.raytracer.Scene == nil {
		return nil, fmt.Errorf("N/A scene")
	}
	image := hdrimage.New(settings.Width, settings.Height)
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

func (cr *ConcurrentRaytracer) SetScene(scene *scene.Scene) {
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
func (cr *ConcurrentRaytracer) GetImage() *hdrimage.Image {
	var mergedSamples *hdrimage.Image

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

	if mergedSamples == nil {
		return hdrimage.New(0, 0)
	}
	return mergedSamples
}

func (cr *ConcurrentRaytracer) ParallelSamples() int {
	return cr.parallelSamples
}
