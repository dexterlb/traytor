package sampler

import (
	"encoding/json"
	"fmt"

	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/jsonutil"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// AnySampler implements the Sampler interface and is deserialiseable from json
type AnySampler struct {
	Sampler
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *AnySampler) UnmarshalJSON(data []byte) error {
	numberSampler := &NumberSampler{}
	err := json.Unmarshal(data, &numberSampler)
	if err == nil {
		*s = AnySampler{numberSampler}
		return nil
	}

	vec3Sampler := &Vec3Sampler{}
	err = json.Unmarshal(data, &vec3Sampler)
	if err == nil {
		*s = AnySampler{vec3Sampler}
		return nil
	}

	samplerType, err := jsonutil.ObjectType(data)
	if err != nil {
		return err
	}

	switch samplerType {
	case "image_texture":
		sampler := &ImageTexture{}
		err = json.Unmarshal(data, &sampler)
		if err != nil {
			return err
		}
		*s = AnySampler{sampler}
	default:
		return fmt.Errorf("Unknown texture sampler: '%s'", samplerType)
	}

	return nil
}

// Sampler objects return varying attributes on points on surfaces
// (e.g. texture colours, ray information etc)
type Sampler interface {
	GetColour(intersection *ray.Intersection) *hdrcolour.Colour
	GetVec3(intersection *ray.Intersection) *maths.Vec3
	GetFac(intersectuib *ray.Intersection) float64
}
