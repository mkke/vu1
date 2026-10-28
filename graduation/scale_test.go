package graduation

import (
	"errors"
	"math"
	"testing"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func TestScaleGraduationFraction(t *testing.T) {
	cases := []struct {
		name  string
		scale ScaleGraduation
		v     float64
		want  float64
	}{
		{"linear middle", ScaleGraduation{Min: 0, Max: 100, Step: 10}, 50, 0.5},
		{"linear below", ScaleGraduation{Min: 0, Max: 100, Step: 10}, -5, 0},
		{"linear above", ScaleGraduation{Min: 0, Max: 100, Step: 10}, 150, 1},
		// 0..55 in steps of 10 ends its arc at the tick for 60.
		{"linear arc past max", ScaleGraduation{Min: 0, Max: 55, Step: 10}, 30, 0.5},
		{"log one decade in", ScaleGraduation{Min: 1, Max: 1000, Log: true}, 10, 1.0 / 3},
		{"log top", ScaleGraduation{Min: 1, Max: 1000, Log: true}, 1000, 1},
		{"log not positive", ScaleGraduation{Min: 1, Max: 1000, Log: true}, -1, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.scale.Fraction(tc.v); math.Abs(got-tc.want) > 1e-12 {
				t.Errorf("Fraction(%g) = %g, want %g", tc.v, got, tc.want)
			}
		})
	}
}

func TestScaleGraduationDrawRefusesInvalidScales(t *testing.T) {
	cases := []struct {
		name  string
		scale ScaleGraduation
	}{
		{"max not above min", ScaleGraduation{Min: 10, Max: 10, Step: 1}},
		{"linear without step", ScaleGraduation{Min: 0, Max: 10}},
		{"log from zero", ScaleGraduation{Min: 0, Max: 100, Log: true}},
		{"log over part of a decade", ScaleGraduation{Min: 1, Max: 50, Log: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.scale.Draw(gg.NewContext(10, 10)); !errors.Is(err, ErrInvalidScale) {
				t.Errorf("Draw = %v, want %v", err, ErrInvalidScale)
			}
		})
	}
}

func TestScaleGraduationDraws(t *testing.T) {
	face := truetype.NewFace(testFont, &truetype.Options{Size: 11})
	base := ScaleGraduation{
		Graduation: VU1, FontFaceL: face, FontFaceM: face,
		Height: 100, LenL: 10, LenM: 5, LenS: 4,
	}

	linear := base
	linear.Min, linear.Max, linear.Step = 0, 100, 10
	linear.Ticks = TickIntervals{IntervalL: 5, IntervalM: 1}
	linear.Marked = []Range{{From: 80, To: 100}}
	linear.Filled = []Range{{From: 95, To: 100}}

	log := base
	log.Min, log.Max, log.Log = 1, 1000, true

	for name, s := range map[string]ScaleGraduation{"linear": linear, "log": log} {
		t.Run(name, func(t *testing.T) {
			dc := gg.NewContext(testWidth, testHeight)
			if err := s.Draw(dc); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFormatTick(t *testing.T) {
	for v, want := range map[float64]string{15: "15", 17.5: "17.5", 0.1: "0.1", 1000: "1000"} {
		if got := formatTick(v); got != want {
			t.Errorf("formatTick(%g) = %q, want %q", v, got, want)
		}
	}
}
