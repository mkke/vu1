package graduation

import "testing"

func TestTickIntervals_TypeForIndex(t *testing.T) {
	ti := TickIntervals{IntervalL: 5, IntervalM: 1}
	cases := []struct {
		idx  int
		want TickType
	}{
		{0, TickL},  // 0 % 5 == 0 → L
		{1, TickM},  // 1 % 5 != 0, 1 % 1 == 0 → M (every non-L is M when IntervalM=1)
		{5, TickL},
		{10, TickL},
	}
	for _, tc := range cases {
		if got := ti.TypeForIndex(tc.idx); got != tc.want {
			t.Errorf("TypeForIndex(%d) = %q, want %q", tc.idx, got, tc.want)
		}
	}
}

func TestTickIntervals_OnlyM(t *testing.T) {
	ti := TickIntervals{IntervalM: 3}
	if ti.TypeForIndex(3) != TickM {
		t.Error("3 should be TickM with IntervalM=3 and IntervalL=0")
	}
	if ti.TypeForIndex(4) != TickS {
		t.Error("4 should fall through to TickS")
	}
}

func TestTickIntervals_AllZeroIsAllSmall(t *testing.T) {
	ti := TickIntervals{}
	for i := 0; i < 10; i++ {
		if ti.TypeForIndex(i) != TickS {
			t.Errorf("with all zero intervals, index %d should be TickS", i)
		}
	}
}
