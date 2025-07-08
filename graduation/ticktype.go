package graduation

type TickType string

const (
	TickL TickType = "L"
	TickM TickType = "M"
	TickS TickType = "S"
)

type TickIntervals struct {
	IntervalL int
	IntervalM int
}

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
