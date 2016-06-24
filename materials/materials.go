package materials

import (
	"encoding/json"
	"fmt"

	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/jsonutil"
	"github.com/DexterLB/traytor/random"
	"github.com/DexterLB/traytor/ray"
)

// Raytracer is any possible raytracer
type Raytracer interface {
	Raytrace(incoming *ray.Ray) *hdrcolour.Colour
	RandomGen() *random.Random
}

// Material objects are used to shade surfaces
type Material interface {
	Shade(intersection *ray.Intersection, raytracer Raytracer) *hdrcolour.Colour
}

// AnyMaterial implements the Material interface and is deserialiseable from json
type AnyMaterial struct {
	Material
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *AnyMaterial) UnmarshalJSON(data []byte) error {
	materialType, err := jsonutil.ObjectType(data)
	if err != nil {
		return err
	}

	switch materialType {
	case "emissive":
		material := &EmissiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "lambert":
		material := &LambertMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "reflective":
		material := &ReflectiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "refractive":
		material := &RefractiveMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	case "mixed":
		material := &MixedMaterial{}
		err = json.Unmarshal(data, &material)
		if err != nil {
			return err
		}
		*m = AnyMaterial{material}
	default:
		return fmt.Errorf("Unknown material type: '%s'", materialType)
	}

	return nil
}
