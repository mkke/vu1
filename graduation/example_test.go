package graduation_test

import (
	"fmt"

	"github.com/mkke/vu1/graduation"
)

func ExampleTickIntervals_TypeForIndex() {
	ti := graduation.TickIntervals{IntervalL: 4, IntervalM: 2}
	for idx := range 5 {
		fmt.Print(ti.TypeForIndex(idx), " ")
	}
	fmt.Println()
	// Output: L S M S L
}
