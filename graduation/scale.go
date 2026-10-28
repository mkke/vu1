package graduation

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// ErrInvalidScale is returned by ScaleGraduation.Draw for a scale it
// cannot draw.
var ErrInvalidScale = errors.New("invalid scale")

// Mark is a stretch of a scale from Min to Max, in the scale's own unit,
// drawn as a band above the arc. A Filled mark is drawn filled, the others
// as an outline.
type Mark struct {
	Min, Max float64
	Filled   bool
}

// ScaleGraduation is a scale from Min to Max, linear or logarithmic.
//
// A linear scale has a tick every Step, starting at Min. The arc ends at
// the first tick at or above Max. Ticks classifies the ticks by their
// index. Large ticks are labelled with FontFaceL, and when LabelMinor is
// set every other tick is labelled with FontFaceM.
//
// A logarithmic scale spans whole decades, so Max must be Min times a
// power of ten. It draws a tick at 1 to 10 times the start of each decade
// and labels each decade boundary with FontFaceL. Step, Ticks and
// LabelMinor do not apply.
//
// Height is the arc's distance from the pivot, and LenL, LenM and LenS
// are the tick lengths in pixels. Each of Marks is drawn as a sector LenM
// deep above the arc, in list order.
type ScaleGraduation struct {
	Graduation
	FontFaceL, FontFaceM font.Face
	Height               int
	Min, Max             float64
	Log                  bool
	Step                 float64
	Ticks                TickIntervals
	LabelMinor           bool
	LenL, LenM, LenS     int
	Marks                []Mark
}

// Fraction returns where v sits on the scale: 0 at Min, 1 at the end of
// the arc. A value outside the scale is clamped to its nearer end.
func (s ScaleGraduation) Fraction(v float64) float64 {
	var f float64
	if s.Log {
		if v <= 0 {
			return 0
		}
		f = math.Log10(v/s.Min) / math.Log10(s.Max/s.Min)
	} else {
		f = (v - s.Min) / (s.top() - s.Min)
	}
	return min(max(f, 0), 1)
}

// top is the value at the end of a linear scale's arc.
func (s ScaleGraduation) top() float64 {
	return s.Min + float64(s.tickCount()-1)*s.Step
}

func (s ScaleGraduation) tickCount() int {
	return int(math.Ceil((s.Max-s.Min)/s.Step)) + 1
}

func (s ScaleGraduation) decades() int {
	return int(math.Round(math.Log10(s.Max / s.Min)))
}

func (s ScaleGraduation) validate() error {
	switch {
	case !(s.Max > s.Min):
		return fmt.Errorf("%w: max %g is not above min %g", ErrInvalidScale, s.Max, s.Min)
	case s.Log && s.Min <= 0:
		return fmt.Errorf("%w: a log scale needs a min above 0, got %g", ErrInvalidScale, s.Min)
	case s.Log && math.Abs(math.Log10(s.Max/s.Min)-float64(s.decades())) > 1e-9:
		return fmt.Errorf("%w: a log scale spans whole decades, %g to %g does not", ErrInvalidScale, s.Min, s.Max)
	case !s.Log && !(s.Step > 0):
		return fmt.Errorf("%w: a linear scale needs a step above 0, got %g", ErrInvalidScale, s.Step)
	}
	return nil
}

// Draw draws the arc, the marks, the ticks and the labels.
func (s ScaleGraduation) Draw(dc *gg.Context) error {
	if err := s.validate(); err != nil {
		return err
	}

	points := s.tickCount()
	if s.Log {
		points = s.decades() * 10
	}
	s.FullArc(dc, s.Height, points)
	for _, m := range s.Marks {
		draw := s.SectorRange
		if m.Filled {
			draw = s.FilledSectorRange
		}
		draw(dc, s.Height, points, s.LenM, s.Fraction(m.Min), s.Fraction(m.Max))
	}

	labelHeight := int(float64(s.Height) * 1.2)
	if s.Log {
		s.LogFullTicks(dc, s.Height, s.decades(), s.LenL, s.LenM, s.LenS)
		s.FullLabel(dc, labelHeight, s.decades()+1, func(dc *gg.Context, idx int) {
			dc.SetFontFace(s.FontFaceL)
			dc.DrawStringAnchored(formatTick(s.Min*math.Pow(10, float64(idx))), 0, 0, 0.5, 0.5)
		})
		return nil
	}

	s.FullTicks(dc, s.Height, points, s.LenL, s.LenM, s.LenS, s.Ticks.IntervalL, s.Ticks.IntervalM)
	s.FullLabel(dc, labelHeight, points, func(dc *gg.Context, idx int) {
		switch tickType := s.Ticks.TypeForIndex(idx); {
		case tickType == TickL:
			dc.SetFontFace(s.FontFaceL)
		case s.LabelMinor:
			dc.SetFontFace(s.FontFaceM)
		default:
			return
		}
		dc.DrawStringAnchored(formatTick(s.Min+float64(idx)*s.Step), 0, 0, 0.5, 0.5)
	})
	return nil
}

// formatTick writes a tick value without trailing zeros.
func formatTick(v float64) string {
	text := strings.TrimRight(fmt.Sprintf("%f", v), "0")
	text, _ = strings.CutSuffix(text, ".")
	return text
}
