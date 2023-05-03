package requiring

import (
	"fmt"
	"io"
	"testing"

	"github.com/lufia/go-pointer"
)

type testInRangeViolationPrinter[T ordered] struct{}

func (testInRangeViolationPrinter[T]) Print(w io.Writer, e InRangeViolationError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

var _ InRangeViolationPrinter[int] = (*testInRangeViolationPrinter[int])(nil)

func TestNotEmptyVerified_int(t *testing.T) {
	tests := []struct {
		value    int
		min, max *int
	}{
		{1, pointer.Int(0), pointer.Int(3)},
		{1, pointer.Int(1), pointer.Int(1)},
	}
	for _, tt := range tests {
		err := InRange(tt.min, tt.max).Validate(tt.value)
		if err != nil {
			t.Fatalf("(%d ... %d).Validate(%d): %v",
				pointer.NewIntFormatter(tt.min),
				pointer.NewIntFormatter(tt.max),
				tt.value, err)
		}
	}
}
