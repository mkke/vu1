package graduation

import (
	"image/color"

	"github.com/fogleman/gg"
)

// VU1 is the geometry and line style for a VU-1 dial.
var VU1 = Graduation{
	MechAngle:     gg.Radians(75),
	MechShiftY:    -.6,
	ArcColor:      color.Black,
	ArcWidth:      2,
	TickColor:     color.Black,
	TickWidthL:    3,
	TickWidthM:    2.5,
	TickWidthLogM: 2.0,
	TickWidthS:    2.0,
	LabelColor:    color.Black,
}
