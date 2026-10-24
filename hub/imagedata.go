package hub

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// Size of a dial's display in pixels.
const (
	DialImageWidth  = 200
	DialImageHeight = 144
)

// ErrInvalidImageDimensions is returned for an image that is not exactly
// DialImageWidth by DialImageHeight pixels.
var ErrInvalidImageDimensions = fmt.Errorf("invalid image dimensions, expected %dx%d", DialImageWidth, DialImageHeight)

// ConvertImageData packs sourceImage into the display's 1-bit format,
// described in the package documentation.
func ConvertImageData(sourceImage image.Image) ([]byte, error) {
	bounds := sourceImage.Bounds()
	if bounds.Dx() != DialImageWidth || bounds.Dy() != DialImageHeight {
		return nil, ErrInvalidImageDimensions
	}

	bwSourceImage := &ConvertedImage{
		Img: sourceImage,
		Mod: color.Palette{color.Black, color.White},
	}

	// Each byte is 8 vertical bits
	var data []byte
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y += 8 {
			d := byte(0)
			for bit := range 8 {
				setBit := func() {
					d |= 128 >> uint(bit)
				}
				pixel := bwSourceImage.At(x, y+bit)
				switch pixel := pixel.(type) {
				case color.Gray:
					if pixel.Y > 127 {
						setBit()
					}
				case color.Gray16:
					if pixel.Y > 32767 {
						setBit()
					}
				}
			}
			data = append(data, d)
		}
	}

	return data, nil
}

// see https://stackoverflow.com/questions/8697095/how-to-read-a-png-file-in-color-and-output-as-gray-scale-using-the-go-programmin

// ConvertedImage implements image.Image, so you can
// pretend that it is the converted image.
type ConvertedImage struct {
	Img image.Image
	Mod color.Model
}

// ColorModel returns the new color model.
func (c *ConvertedImage) ColorModel() color.Model {
	return c.Mod
}

// Bounds returns the original bounds
func (c *ConvertedImage) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image and
// then asks the color model to convert it.
func (c *ConvertedImage) At(x, y int) color.Color {
	return c.Mod.Convert(c.Img.At(x, y))
}

// DrawImage draws an image with dial dimensions using the given draw function
func DrawImage(draw func(dc *gg.Context) error) (image.Image, error) {
	dc := gg.NewContext(DialImageWidth, DialImageHeight)
	dc.SetColor(color.White)
	dc.Clear()

	dc.SetColor(color.Black)
	if err := draw(dc); err != nil {
		return nil, err
	}

	return dc.Image(), nil
}
