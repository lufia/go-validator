package validator

import (
	"fmt"
	"io"
	"testing"
)

func TestMin(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Min(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 4, "")
		testValidate(t, v, 2, "must be no less than 3")
	})
}

type testMinErrorPrinter[T ordered] struct{}

func (testMinErrorPrinter[T]) Print(w io.Writer, e *MinError[T]) {
	fmt.Fprintf(w, "'%v' is less than %v", e.Value, e.Min)
}

func TestMinWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := Min(3).WithPrinter(&testMinErrorPrinter[int]{})
		testValidate(t, v, 2, "'2' is less than 3")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Min(3).WithPrinterFunc(func(w io.Writer, min int) {
			fmt.Fprintf(w, "less than %d", min)
		})
		testValidate(t, v, 2, "less than 3")
	})
}

func TestMax(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Max(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 2, "")
		testValidate(t, v, 4, "must be no greater than 3")
	})
}

type testMaxErrorPrinter[T ordered] struct{}

func (testMaxErrorPrinter[T]) Print(w io.Writer, e *MaxError[T]) {
	fmt.Fprintf(w, "'%v' is greater than %v", e.Value, e.Max)
}

func TestMaxWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := Max(3).WithPrinter(&testMaxErrorPrinter[int]{})
		testValidate(t, v, 4, "'4' is greater than 3")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Max(3).WithPrinterFunc(func(w io.Writer, max int) {
			fmt.Fprintf(w, "greater than %d", max)
		})
		testValidate(t, v, 4, "greater than 3")
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
}

type testInRangeErrorPrinter[T ordered] struct{}

func (testInRangeErrorPrinter[T]) Print(w io.Writer, e *InRangeError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

func TestInRangeWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := InRange(1, 3).WithPrinter(&testInRangeErrorPrinter[int]{})
		testValidate(t, v, 4, "'4' is out of range(1, 3)")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := InRange(1, 3).WithPrinterFunc(func(w io.Writer, min, max int) {
			fmt.Fprintf(w, "out of range(%d, %d)", min, max)
		})
		testValidate(t, v, 4, "out of range(1, 3)")
	})
}
