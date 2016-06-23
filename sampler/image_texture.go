package sampler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"strings"

	"github.com/DexterLB/traytor/hdrcolour"
	"github.com/DexterLB/traytor/hdrimage"
	"github.com/DexterLB/traytor/maths"
	"github.com/DexterLB/traytor/ray"
)

// ImageTexture is.. an image texture!
type ImageTexture struct {
	Image                            *hdrimage.Image
	ScaleU, ScaleV, OffsetU, OffsetV float64
}

// GetColour returns the colour at (u, v)
func (i *ImageTexture) GetColour(intersection *ray.Intersection) *hdrcolour.Colour {
	u := i.OffsetU + i.ScaleU*intersection.U
	v := i.OffsetV + i.ScaleV*intersection.V

	// tile it
	u -= math.Floor(u)
	v -= math.Floor(v)

	if u < 0 {
		u++
	}
	if v < 0 {
		v++
	}

	v = 1 - v // (0, 0) is at the topleft corner of images

	return i.Image.AtHDR(int(u*float64(i.Image.Width-1)), int(v*float64(i.Image.Height-1)))
}

// GetFac returns the colour intensity at (u, v)
func (i *ImageTexture) GetFac(intersection *ray.Intersection) float64 {
	return float64(i.GetColour(intersection).Intensity())
}

// GetVec3 returns the colour at (u, v)'s components as a vector
func (i *ImageTexture) GetVec3(intersection *ray.Intersection) *maths.Vec3 {
	colour := i.GetColour(intersection)
	return maths.NewVec3(float64(colour.R), float64(colour.G), float64(colour.B))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *ImageTexture) UnmarshalJSON(data []byte) error {
	textureSettings := &struct {
		Scale  *maths.Vec3
		Offset *maths.Vec3
		Format string
		Data   []byte
	}{}

	err := json.Unmarshal(data, textureSettings)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(textureSettings.Data)

	var ldrImage image.Image

	switch strings.ToLower(textureSettings.Format) {
	case "png":
		ldrImage, err = png.Decode(reader)
		if err != nil {
			return err
		}
	case "jpeg":
		ldrImage, err = jpeg.Decode(reader)
		if err != nil {
			return err
		}
	case "traytor_hdr":
		i.Image, err = hdrimage.Decode(reader)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"Unknown format for image texture: %s\n",
			textureSettings.Format,
		)
	}

	if ldrImage != nil {
		i.Image = hdrimage.FromSRGB(ldrImage)
	}

	i.ScaleU = textureSettings.Scale.X
	i.ScaleV = textureSettings.Scale.Y

	i.OffsetU = textureSettings.Offset.X
	i.OffsetV = textureSettings.Offset.Y

	return nil
}
