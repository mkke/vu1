package graduation

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

// Graduation describes a dial's needle geometry and the styles a scale is
// drawn with.
type Graduation struct {
	// MechAngle is the needle's full sweep in radians.
	MechAngle float64
	// MechShiftY is the vertical position of the needle pivot as a fraction
	// of the image height, measured up from the bottom edge.
	MechShiftY float64

	ArcColor color.Color
	ArcWidth float64

	TickColor     color.Color
	TickWidthL    float64
	TickWidthM    float64
	TickWidthLogM float64
	TickWidthS    float64

	LabelColor color.Color
}

func (g *Graduation) pushCoordinateSystem(dc *gg.Context) {
	dc.Push()
	dc.InvertY()
	dc.Translate(float64(dc.Width())/2, g.MechShiftY*float64(dc.Height()))
}

// ToXY converts polar to cartesian.
//   - height is the height at which the point is required
//   - angle is between 0 (min deflection) and 1 (max deflection)
func (g *Graduation) ToXY(height, angle float64) (x, y float64) {
	a := angle*g.MechAngle - g.MechAngle/2.0
	y = math.Cos(a) * height
	x = math.Sin(a) * height
	return x, y
}

// ToMechAngle converts a scale position from 0 to 1 to the needle angle in
// radians, 0 being straight up.
func (g *Graduation) ToMechAngle(angle float64) float64 {
	return angle*g.MechAngle - g.MechAngle/2.0
}

func (g *Graduation) drawPartArc(dc *gg.Context, height, points int, segmentIdx, numSegments int, clockwise bool) {
	from := float64(segmentIdx) / float64(numSegments)
	to := float64(segmentIdx+1) / float64(numSegments)
	g.drawArcRange(dc, height, points, from, to, clockwise)
}

// drawArcRange adds an arc at height from scale position from to scale
// position to, both between 0 and 1, to the current path.
func (g *Graduation) drawArcRange(dc *gg.Context, height, points int, from, to float64, clockwise bool) {
	drawPoint := func(idx, points int) {
		x, y := g.ToXY(float64(height), from+(to-from)*float64(idx)/float64(points-1))

		if _, hasCurrent := dc.GetCurrentPoint(); !hasCurrent {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}

	points *= 10
	if clockwise {
		for idx := range points {
			drawPoint(idx, points)
		}
	} else {
		for idx := points - 1; idx >= 0; idx-- {
			drawPoint(idx, points)
		}
	}
}

// FullArc draws a full arc
func (g *Graduation) FullArc(dc *gg.Context, height, points int) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.ArcColor)
	dc.SetLineWidth(g.ArcWidth)
	dc.SetLineCapRound()

	g.drawPartArc(dc, height, points, 0, 1, true)

	dc.Stroke()
	dc.Pop()
}

// PartArc draws a partial arc at height h, with p points
func (g *Graduation) PartArc(dc *gg.Context, height, points int, segmentIdx, numSegments int) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.ArcColor)
	dc.SetLineWidth(g.ArcWidth)
	dc.SetLineCapRound()

	g.drawPartArc(dc, height, points, segmentIdx, numSegments, true)

	dc.Stroke()
	dc.Pop()
}

// Sector draws a sector
func (g *Graduation) Sector(dc *gg.Context, height, points, lift int, segmentIdx, numSegments int) {
	from := float64(segmentIdx) / float64(numSegments)
	to := float64(segmentIdx+1) / float64(numSegments)
	g.SectorRange(dc, height, points, lift, from, to)
}

// SectorRange draws the outline of a band lift pixels deep above the arc
// at height, from scale position from to scale position to, both between
// 0 and 1.
func (g *Graduation) SectorRange(dc *gg.Context, height, points, lift int, from, to float64) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.ArcColor)
	dc.SetLineWidth(g.ArcWidth)
	dc.SetLineCapRound()

	// bottom arc
	g.drawArcRange(dc, height, points, from, to, true)
	// top arc
	g.drawArcRange(dc, height+lift, points, from, to, false)
	dc.ClosePath()

	dc.Stroke()
	dc.Pop()
}

// FullTicks draws major, semi and minor ticks
func (g *Graduation) FullTicks(dc *gg.Context, height, points, lenL, lenM, lenS, intervalL, intervalM int) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.TickColor)
	dc.SetLineCapRound()

	ti := TickIntervals{IntervalL: intervalL, IntervalM: intervalM}

	for idx := range points {
		angle := float64(idx) / float64(points-1)
		x1, y1 := g.ToXY(float64(height), angle)

		var x2, y2 float64
		switch ti.TypeForIndex(idx) {
		case TickL:
			dc.SetLineWidth(g.TickWidthL)
			x2, y2 = g.ToXY(float64(height+lenL), angle)
		case TickM:
			dc.SetLineWidth(g.TickWidthM)
			x2, y2 = g.ToXY(float64(height+lenM), angle)
		default:
			dc.SetLineWidth(g.TickWidthS)
			x2, y2 = g.ToXY(float64(height+lenS), angle)
		}
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	dc.Pop()
}

// LogFullTicks draws major, semi and minor ticks for a log graduation
func (g *Graduation) LogFullTicks(dc *gg.Context, height, decades, lenL, lenM, lenS int) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.TickColor)

	points := decades * 10
	dsel := float64(0)
	inc := float64(1)

	ti := TickIntervals{IntervalL: 10, IntervalM: 5}

	for idx := range points {
		angle := (dsel / float64(decades)) + math.Log10(inc)*(1.0/float64(decades))
		x1, y1 := g.ToXY(float64(height), angle)

		var x2, y2 float64
		switch ti.TypeForIndex(idx) {
		case TickL:
			dc.SetLineWidth(g.TickWidthL)
			x2, y2 = g.ToXY(float64(height+lenL), angle)
		case TickM:
			dc.SetLineWidth(g.TickWidthM)
			x2, y2 = g.ToXY(float64(height+lenM), angle)
		default:
			dc.SetLineWidth(g.TickWidthS)
			x2, y2 = g.ToXY(float64(height+lenS), angle)
		}
		inc++
		if inc > 10 {
			inc = 1
			dsel++
		}
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	dc.Pop()
}

// FullLabel calls drawLabelFn once per point, with the context moved to the
// point at height and rotated to the needle angle there.
func (g *Graduation) FullLabel(dc *gg.Context, height, points int, drawLabelFn func(dc *gg.Context, idx int)) {
	g.pushCoordinateSystem(dc)
	dc.SetColor(g.LabelColor)

	for idx := range points {
		x, y := g.ToXY(float64(height), float64(idx)/float64(points-1))
		angle := float64(idx) / float64(points-1)
		mechAngle := g.ToMechAngle(angle)

		dc.Push()
		dc.Translate(x, y)
		dc.Scale(1, -1)
		dc.Rotate(mechAngle)

		drawLabelFn(dc, idx)
		dc.Pop()
	}

	dc.Pop()
}
