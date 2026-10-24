package graduation

// TickType is the size class of a tick: large, medium or small.
type TickType string

// Tick size classes.
const (
	TickL TickType = "L"
	TickM TickType = "M"
	TickS TickType = "S"
)

// TickIntervals gives every IntervalL-th point a large tick and every
// IntervalM-th of the rest a medium one. An interval of 0 disables that
// class.
type TickIntervals struct {
	IntervalL int
	IntervalM int
}

// TypeForIndex returns the tick class of point idx.
func (ti TickIntervals) TypeForIndex(idx int) TickType {
	switch {
	case ti.IntervalL > 0 && idx%ti.IntervalL == 0:
		return TickL
	case ti.IntervalM > 0 && idx%ti.IntervalM == 0:
		return TickM
	default:
		return TickS
	}
}
