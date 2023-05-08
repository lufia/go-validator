package validator

import (
	"fmt"
	"io"
	"testing"
)

func TestRequired(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Required[string]()
		testValidate(t, v, "a", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "", "cannot be the zero value")
	})
	t.Run("int", func(t *testing.T) {
		v := Required[int]()
		testValidate(t, v, 1, "")
		testValidate(t, v, -1, "")
	})
}

type testRequiredErrorPrinter[T comparable] struct{}

func (testRequiredErrorPrinter[T]) Print(w io.Writer, e *RequiredError[T]) {
	fmt.Fprintf(w, "'%v' is empty", e.Value)
}

func TestRequiredWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := Required[string]().WithPrinter(&testRequiredErrorPrinter[string]{})
		testValidate(t, v, "", "'' is empty")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Required[string]().WithPrinterFunc(func(w io.Writer) {
			fmt.Fprintf(w, "is empty")
		})
		testValidate(t, v, "", "is empty")
	})
}
