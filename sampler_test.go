package traytor

import (
	"encoding/json"
	"fmt"
)

func ExampleVec3Sampler() {
	var sampler *Vec3Sampler
	err := json.Unmarshal([]byte(`[0.1, 0.5, 0.3]`), &sampler)
	if err != nil {
		fmt.Printf("can't unmarshal data: %s\n", err)
		return
	}

	fmt.Printf("Vector sample: %s\n", sampler.GetVec3(nil))
	fmt.Printf("Colour sample: %s\n", sampler.GetColour(nil))
	fmt.Printf("Fac sample: %.3g\n", sampler.GetFac(nil))

	// Output:
	// Vector sample: (0.1, 0.5, 0.3)
	// Colour sample: {0.1, 0.5, 0.3}
	// Fac sample: 0.3
}

func ExampleNumberSampler() {
	var sampler *NumberSampler
	err := json.Unmarshal([]byte(`42`), &sampler)
	if err != nil {
		fmt.Printf("can't unmarshal data: %s\n", err)
		return
	}

	fmt.Printf("Vector sample: %s\n", sampler.GetVec3(nil))
	fmt.Printf("Colour sample: %s\n", sampler.GetColour(nil))
	fmt.Printf("Fac sample: %.3g\n", sampler.GetFac(nil))

	// Output:
	// Vector sample: (42, 42, 42)
	// Colour sample: {42, 42, 42}
	// Fac sample: 42
}
