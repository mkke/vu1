package graduation

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/mkke/vu1/hub"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	testScale  = 1
	testWidth  = hub.DialImageWidth * testScale
	testHeight = hub.DialImageHeight * testScale
)

var testFont, _ = truetype.Parse(goregular.TTF)

func drawPNG(path string, label string, draw func(dc *gg.Context) error) error {
	dc := gg.NewContext(testWidth, testHeight)
	dc.SetColor(color.White)
	dc.Clear()

	dc.SetFontFace(truetype.NewFace(testFont, &truetype.Options{Size: 14 * testScale}))

	dc.SetColor(color.Black)
	dc.DrawStringAnchored(label, float64(testWidth)/2, float64(testHeight)*0.2, 0.5, 0.5)

	if err := draw(dc); err != nil {
		return err
	}

	img := dc.Image()

	if f, err := os.Create(path); err != nil {
		return err
	} else {
		defer f.Close()
		if err := png.Encode(f, img); err != nil {
			return err
		}
	}
	return nil
}

func drawBasicAngle(dc *gg.Context) {
	VU1.pushCoordinateSystem(dc)
	defer dc.Pop()

	height := float64(testHeight) * 0.6
	drawAngle := func(angle float64, lineWidth float64) {
		x1, y1 := VU1.ToXY(height, angle)
		x2, y2 := VU1.ToXY(height*1.2, angle)
		dc.SetLineWidth(lineWidth)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}
	drawAngle(0, 1)
	drawAngle(0.5, 2)
	drawAngle(1, 1)
}

func TestToXY(t *testing.T) {
	if err := drawPNG("TestToXY.png", "TestToXY", func(dc *gg.Context) error {
		drawBasicAngle(dc)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGraduation(t *testing.T) {
	if err := drawPNG("TestGraduation.png", "TestGraduation", func(dc *gg.Context) error {
		height := float64(testHeight) * 0.8

		TemperatureGraduation{
			Graduation: VU1,
			FontFaceL:  truetype.NewFace(testFont, &truetype.Options{Size: 14 * testScale}),
			FontFaceM:  truetype.NewFace(testFont, &truetype.Options{Size: 11 * testScale}),
			Height:     int(height),
			LowTemp:    15,
			HighTemp:   35,
			Step:       2.5,
			LenL:       10,
			LenM:       5,
			LenS:       4,
			Sectors:    nil,
		}.Draw(dc)

		dc.SetFontFace(truetype.NewFace(testFont, &truetype.Options{Size: 18 * testScale}))
		dc.SetColor(color.Black)
		dc.DrawStringAnchored("°C", float64(testWidth)/2, float64(testHeight)*0.8, 0.5, 0.5)

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
