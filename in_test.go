package validator

import (
	"fmt"
	"io"
	"testing"
)

func TestIn(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		v := In("a", "b")
		testValidate(t, v, "a", "")
		testValidate(t, v, "b", "")
		testValidate(t, v, "x", "must be a valid value in [a b]")
	})
	t.Run("int", func(t *testing.T) {
		v := In(1, 2)
		testValidate(t, v, 1, "")
		testValidate(t, v, 2, "")
	})
}

type testInErrorPrinter[T comparable] struct{}

func (testInErrorPrinter[T]) Print(w io.Writer, e *InError[T]) {
	fmt.Fprintf(w, "'%v' is not valid", e.Value)
}

func TestInWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := In("a", "b").WithPrinter(&testInErrorPrinter[string]{})
		testValidate(t, v, "x", "'x' is not valid")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := In("a", "b").WithPrinterFunc(func(w io.Writer, a []string) {
			fmt.Fprintf(w, "is not in %v", a)
		})
		testValidate(t, v, "", "is not in [a b]")
	})
}
