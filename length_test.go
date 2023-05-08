package validator

import (
	"fmt"
	"io"
	"testing"
)

func TestMinLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MinLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "1234", "")
		testValidate(t, v, "ab", "the length must be no less than 3")
	})
}

type testMinLengthErrorPrinter[T ~string] struct{}

func (testMinLengthErrorPrinter[T]) Print(w io.Writer, e *MinLengthError[T]) {
	fmt.Fprintf(w, "'%v' is less than %v", e.Value, e.Min)
}

func TestMinLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := MinLength[string](3).WithPrinter(&testMinLengthErrorPrinter[string]{})
		testValidate(t, v, "ab", "'ab' is less than 3")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := MinLength[string](3).WithPrinterFunc(func(w io.Writer, min int) {
			fmt.Fprintf(w, "less than %d", min)
		})
		testValidate(t, v, "ab", "less than 3")
	})
}

func TestMaxLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MaxLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "1234", "the length must be no greater than 3")
	})
}

type testMaxLengthErrorPrinter[T ~string] struct{}

func (testMaxLengthErrorPrinter[T]) Print(w io.Writer, e *MaxLengthError[T]) {
	fmt.Fprintf(w, "'%v' is greater than %v", e.Value, e.Max)
}

func TestMaxLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := MaxLength[string](3).WithPrinter(&testMaxLengthErrorPrinter[string]{})
		testValidate(t, v, "1234", "'1234' is greater than 3")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := MaxLength[string](3).WithPrinterFunc(func(w io.Writer, max int) {
			fmt.Fprintf(w, "greater than %d", max)
		})
		testValidate(t, v, "1234", "greater than 3")
	})
}

func TestLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Length[string](1, 3)
		testValidate(t, v, "a", "")
		testValidate(t, v, "abc", "")
		testValidate(t, v, "", "the length must be in range(1 ... 3)")
		testValidate(t, v, "1234", "the length must be in range(1 ... 3)")
	})
}

type testLengthErrorPrinter[T ~string] struct{}

func (testLengthErrorPrinter[T]) Print(w io.Writer, e *LengthError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

func TestLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := Length[string](1, 3).WithPrinter(&testLengthErrorPrinter[string]{})
		testValidate(t, v, "1234", "'1234' is out of range(1, 3)")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Length[string](1, 3).WithPrinterFunc(func(w io.Writer, min, max int) {
			fmt.Fprintf(w, "out of range(%d, %d)", min, max)
		})
		testValidate(t, v, "1234", "out of range(1, 3)")
	})
}
