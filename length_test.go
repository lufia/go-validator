package requiring

import (
	"fmt"
	"io"
	"testing"
)

type testMinLengthViolationPrinter[T ~string] struct{}

func (testMinLengthViolationPrinter[T]) Print(w io.Writer, e MinLengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is less than %v", e.Value, e.Min)
}

type testMaxLengthViolationPrinter[T ~string] struct{}

func (testMaxLengthViolationPrinter[T]) Print(w io.Writer, e MaxLengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is greater than %v", e.Value, e.Max)
}

type testLengthViolationPrinter[T ~string] struct{}

func (testLengthViolationPrinter[T]) Print(w io.Writer, e LengthViolationError[T]) {
	fmt.Fprintf(w, "'%v' is out of range(%v, %v)", e.Value, e.Min, e.Max)
}

var (
	_ MinLengthViolationPrinter[string] = (*testMinLengthViolationPrinter[string])(nil)
	_ MaxLengthViolationPrinter[string] = (*testMaxLengthViolationPrinter[string])(nil)
	_ LengthViolationPrinter[string]    = (*testLengthViolationPrinter[string])(nil)
)

func TestMinLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MinLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "1234", "")
		testValidate(t, v, "ab", "the length must be no less than 3")
	})
	t.Run("custom printer", func(t *testing.T) {
		v := MinLength[string](3,
			&testMinLengthViolationPrinter[string]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, "ab", "'ab' is less than 3")
		testValidate(t, v, 123, "'' int(123) vs string")
	})
}

func TestMaxLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MaxLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "1234", "the length must be no greater than 3")
	})
	t.Run("custom printer", func(t *testing.T) {
		v := MaxLength[string](3,
			&testMaxLengthViolationPrinter[string]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, "1234", "'1234' is greater than 3")
		testValidate(t, v, 123, "'' int(123) vs string")
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
	t.Run("custom printer", func(t *testing.T) {
		v := Length[string](1, 3,
			&testLengthViolationPrinter[string]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, "1234", "'1234' is out of range(1, 3)")
		testValidate(t, v, 3, "'' int(3) vs string")
	})
}
