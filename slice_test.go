package validator

import (
	"fmt"
	"io"
	"testing"
)

func TestSlice(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Slice[string](Required[string]())
		testValidate(t, v, []string{"a", "ab"}, "")
		testValidate(t, v, []string(nil), "")
		testValidate(t, v, []string{""}, "cannot be the zero value")
	})
	t.Run("int", func(t *testing.T) {
		v := Slice[int](Required[int]())
		testValidate(t, v, []int{1, -1}, "")
	})
}

type testSliceErrorPrinter[T any] struct{}

func (testSliceErrorPrinter[T]) Print(w io.Writer, e *SliceError[T]) {
	fmt.Fprintf(w, "%v: ", e.Value)
	for _, key := range e.Errors.Keys() {
		v, _ := e.Errors.Get(key)
		fmt.Fprintf(w, "%v; ", v)
	}
	fmt.Fprintln(w)
}

func TestSliceWithPrinter(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		p := &testSliceErrorPrinter[string]{}
		v := Slice[string](Required[string]()).WithPrinter(p)
		testValidate(t, v, []string{""}, "[]: cannot be the zero value; \n")
	})
	t.Run("printerfunc", func(t *testing.T) {
		v := Slice[string](Required[string]()).WithPrinterFunc(func(w io.Writer) {
			fmt.Fprintf(w, "is empty")
		})
		testValidate(t, v, []string{""}, "is empty")
	})
}
