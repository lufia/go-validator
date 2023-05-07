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

type testMinLengthViolationPrinter[T ~string] struct{}

func (testMinLengthViolationPrinter[T]) Print(w io.Writer, e *MinLengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is less than %v", e.Value, e.Min)
}

func TestMinLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := MinLength[string](3).WithPrinter(&testMinLengthViolationPrinter[string]{})
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

type testMaxLengthViolationPrinter[T ~string] struct{}

func (testMaxLengthViolationPrinter[T]) Print(w io.Writer, e *MaxLengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is greater than %v", e.Value, e.Max)
}

func TestMaxLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := MaxLength[string](3).WithPrinter(&testMaxLengthViolationPrinter[string]{})
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

type testLengthViolationPrinter[T ~string] struct{}

func (testLengthViolationPrinter[T]) Print(w io.Writer, e *LengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

func TestLengthWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := Length[string](1, 3).WithPrinter(&testLengthViolationPrinter[string]{})
		testValidate(t, v, "1234", "'1234' is out of range(1, 3)")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Length[string](1, 3).WithPrinterFunc(func(w io.Writer, min, max int) {
			fmt.Fprintf(w, "out of range(%d, %d)", min, max)
		})
		testValidate(t, v, "1234", "out of range(1, 3)")
	})
}
