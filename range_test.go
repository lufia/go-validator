package validator

import (
	"fmt"
	"io"
	"testing"
)

type testMinViolationPrinter[T ordered] struct{}

func (testMinViolationPrinter[T]) Print(w io.Writer, e MinViolationError[T]) {
	fmt.Fprintf(w, "'%v' is less than %v", e.Value, e.Min)
}

type testMaxViolationPrinter[T ordered] struct{}

func (testMaxViolationPrinter[T]) Print(w io.Writer, e MaxViolationError[T]) {
	fmt.Fprintf(w, "'%v' is greater than %v", e.Value, e.Max)
}

type testInRangeViolationPrinter[T ordered] struct{}

func (testInRangeViolationPrinter[T]) Print(w io.Writer, e InRangeViolationError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

var (
	_ MinViolationPrinter[int]     = (*testMinViolationPrinter[int])(nil)
	_ MaxViolationPrinter[int]     = (*testMaxViolationPrinter[int])(nil)
	_ InRangeViolationPrinter[int] = (*testInRangeViolationPrinter[int])(nil)
)

func TestMin(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Min(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 4, "")
		testValidate(t, v, 2, "must be no less than 3")
	})
	t.Run("int with printer", func(t *testing.T) {
		v := Min(3,
			&testMinViolationPrinter[int]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, 2, "'2' is less than 3")
		testValidate(t, v, 3.0, "float64(3) vs int")
	})
}

func TestMax(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Max(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 2, "")
		testValidate(t, v, 4, "must be no greater than 3")
	})
	t.Run("int with printer", func(t *testing.T) {
		v := Max(3,
			&testMaxViolationPrinter[int]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, 4, "'4' is greater than 3")
		testValidate(t, v, 3.0, "float64(3) vs int")
	})
}

func TestInRange(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := InRange(1, 3)
		testValidate(t, v, 1, "")
		testValidate(t, v, 3, "")
		testValidate(t, v, 0, "must be in range(1 ... 3)")
		testValidate(t, v, 4, "must be in range(1 ... 3)")
	})
	t.Run("int with printer", func(t *testing.T) {
		v := InRange(1, 3,
			&testInRangeViolationPrinter[int]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, 4, "'4' is out of range(1, 3)")
		testValidate(t, v, 3.0, "float64(3) vs int")
	})
}
