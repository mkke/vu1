package graduation

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// Sector marks segment SegmentIndex of NumSegments equal segments of the
// scale, such as the last eighth for a warning range.
type Sector struct {
	SegmentIndex int
	NumSegments  int
}

// TemperatureGraduation is a linear scale from LowTemp to HighTemp with a
// tick every Step. Every int(5/Step)-th tick is large and every
// int(1/Step)-th medium. A Step of 2.5 gives a large tick every 5 degrees.
// A Step above 5 gives small ticks only. Large ticks are
// labelled with FontFaceL and the others with FontFaceM. Height is the arc's
// distance from the pivot, and LenL, LenM and LenS are the tick lengths in
// pixels.
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

// Draw draws a temperature graduation through ScaleGraduation, and returns
// its error for a scale it cannot draw.
func (tg TemperatureGraduation) Draw(dc *gg.Context) error {
	s := ScaleGraduation{
		Graduation: tg.Graduation,
		FontFaceL:  tg.FontFaceL,
		FontFaceM:  tg.FontFaceM,
		Height:     tg.Height,
		Min:        tg.LowTemp,
		Max:        tg.HighTemp,
		Step:       tg.Step,
		Ticks:      TickIntervals{IntervalL: int(5.0 / tg.Step), IntervalM: int(1.0 / tg.Step)},
		LabelMinor: true,
		LenL:       tg.LenL,
		LenM:       tg.LenM,
		LenS:       tg.LenS,
	}
	if tg.Step > 0 {
		span := s.top() - s.Min
		for _, sector := range tg.Sectors {
			from := float64(sector.SegmentIndex) / float64(sector.NumSegments)
			to := float64(sector.SegmentIndex+1) / float64(sector.NumSegments)
			s.Marked = append(s.Marked, Range{From: s.Min + from*span, To: s.Min + to*span})
		}
	}
	return s.Draw(dc)
}
