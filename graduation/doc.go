// Package graduation draws needle scales (arcs, ticks and labels) for analog
// dial faces with github.com/fogleman/gg.
//
// A Graduation describes the dial's mechanics and the line styles. VU1 holds
// the values for a VU-1 dial: the needle sweeps 75 degrees, and its pivot
// sits 0.6 times the image height below the bottom edge, centred.
//
// Positions along the scale are fractions from 0 (needle at its minimum) to
// 1 (needle at its maximum). A height is the distance from the pivot in
// pixels. ToXY converts a height and a position to coordinates in the
// Graduation's own coordinate system, which every drawing method sets up on
// the gg.Context and restores before returning.
//
// The drawing methods take a number of points: the count of evenly spaced
// tick positions from 0 to 1, both ends included. TickIntervals decides
// which of them get a large, medium or small tick.
//
// ScaleGraduation combines the pieces into a complete scale from Min to
// Max, linear or logarithmic, with marked ranges drawn as sectors. Its
// Fraction maps a value to the needle position. TemperatureGraduation is
// a ScaleGraduation set up for temperatures from LowTemp to HighTemp:
//
//	img, err := hub.DrawImage(func(dc *gg.Context) error {
//		graduation.TemperatureGraduation{
//			Graduation: graduation.VU1,
//			FontFaceL:  large,
//			FontFaceM:  medium,
//			Height:     140,
//			LowTemp:    15,
//			HighTemp:   35,
//			Step:       2.5,
//			LenL:       12,
//			LenM:       8,
//			LenS:       6,
//		}.Draw(dc)
//		return nil
//	})
package graduation
