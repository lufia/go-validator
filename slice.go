package validator

import (
	"bytes"
	"fmt"
	"io"
)

func Slice[T any](vs ...Validator) *SliceValidator[T] {
	var r SliceValidator[T]
	r.vs = vs
	return &r
}

// SliceValidator represents the validator to check slice elements.
type SliceValidator[T any] struct {
	vs []Validator
	p  SliceErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *SliceValidator[T]) WithPrinter(p SliceErrorPrinter[T]) *SliceValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *SliceValidator[T]) WithPrinterFunc(fn func(w io.Writer)) *SliceValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *SliceError[T]) {
		fn(w)
	})
	return &rr
}

// Validate validates v. If v's type is not []T, Validate panics.
func (r *SliceValidator[T]) Validate(v any) error {
	a := v.([]T)
	var m OrderedMap[int, error]
	for i, elem := range a {
		var errs []error
		for _, rule := range r.vs {
			if err := rule.Validate(elem); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			m.set(i, joinErrors(errs...))
		}
	}
	if m.Len() > 0 {
		return &SliceError[T]{
			Value:  a,
			Errors: &m,
			rule:   r,
		}
	}
	return nil
}

// SliceError reports an error is caused in SliceValidator.
type SliceError[T any] struct {
	Value  []T
	Errors *OrderedMap[int, error]
	rule   *SliceValidator[T]
}

// Error implements the error interface.
func (e SliceError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &sliceErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

// Unwrap returns each errors of err.
func (e SliceError[T]) Unwrap() []error {
	n := e.Errors.Len()
	if n == 0 {
		return nil
	}
	errs := make([]error, 0, n)
	for _, key := range e.Errors.Keys() {
		err, _ := e.Errors.Get(key)
		errs = append(errs, err)
	}
	return errs
}

type sliceErrorPrinter[T any] struct{}

func (sliceErrorPrinter[T]) Print(w io.Writer, e *SliceError[T]) {
	i := 0
	for _, key := range e.Errors.Keys() {
		if i > 0 {
			w.Write([]byte("\n"))
		}
		err, _ := e.Errors.Get(key)
		fmt.Fprint(w, err)
		i++
	}
}

// SliceErrorPrinter is the interface that wraps Print method.
type SliceErrorPrinter[T any] interface {
	Printer[SliceError[T]]
}

var _ typedValidator[
	*SliceValidator[any],
	SliceError[any],
	SliceErrorPrinter[any],
] = (*SliceValidator[any])(nil)
