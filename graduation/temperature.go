package graduation

import (
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type Sector struct {
	SegmentIndex int
	NumSegments  int
}

type TemperatureGraduation struct {
	Graduation
	FontFaceL         font.Face
	FontFaceM         font.Face
	Height            int
	LowTemp, HighTemp float64
	Step              float64
	LenL, LenM, LenS  int
	Sectors           []Sector
}

// Draw draws a temperature graduation.
func (tg TemperatureGraduation) Draw(dc *gg.Context) {
	ticks := int(math.Ceil((tg.HighTemp-tg.LowTemp)/tg.Step)) + 1
	tg.FullArc(dc, tg.Height, ticks)

	for _, sector := range tg.Sectors {
		tg.Sector(dc, tg.Height, ticks, tg.LenM, sector.SegmentIndex, sector.NumSegments)
	}

	ti := TickIntervals{IntervalL: int(5.0 / tg.Step), IntervalM: int(1.0 / tg.Step)}
	tg.FullTicks(dc, tg.Height, ticks, tg.LenL, tg.LenM, tg.LenS, ti.IntervalL, ti.IntervalM)

	tg.FullLabel(dc, int(float64(tg.Height)*1.2), ticks, func(dc *gg.Context, idx int) {
		temp := tg.LowTemp + float64(idx)*tg.Step

		switch ti.TypeForIndex(idx) {
		case TickL:
			dc.SetFontFace(tg.FontFaceL)
		default:
			dc.SetFontFace(tg.FontFaceM)
		}

		text := fmt.Sprintf("%f", temp)
		text = strings.TrimRight(text, "0")
		text, _ = strings.CutSuffix(text, ".")

		dc.DrawStringAnchored(text, 0, 0, 0.5, 0.5)
	})
}
