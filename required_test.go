package validator

import (
	"fmt"
	"io"
	"testing"
)

type testRequiredViolationPrinter[T comparable] struct{}

func (testRequiredViolationPrinter[T]) Print(w io.Writer, e RequiredViolationError[T]) {
	fmt.Fprintf(w, "'%v' is empty", e.Value)
}

var _ RequiredViolationPrinter[string] = (*testRequiredViolationPrinter[string])(nil)

func TestRequired(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Required[string]()
		testValidate(t, v, "a", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "", "cannot be the zero value")
	})
	t.Run("string with printer", func(t *testing.T) {
		v := Required[string](
			&testRequiredViolationPrinter[string]{},
			&testInvalidTypePrinter{},
		)
		testValidate(t, v, "", "'' is empty")
		testValidate(t, v, 0, "int(0) vs string")
	})
	t.Run("int", func(t *testing.T) {
		v := Required[int]()
		testValidate(t, v, 1, "")
		testValidate(t, v, -1, "")
	})
}
