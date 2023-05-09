package validator

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/exp/slices"
)

// In returns the validator to verify the value is in a.
func In[T comparable](a ...T) *InValidator[T] {
	return &InValidator[T]{
		a: a,
	}
}

// InValidator represents the validator to check the value is in T.
type InValidator[T comparable] struct {
	a []T
	p InErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *InValidator[T]) WithPrinter(p InErrorPrinter[T]) *InValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *InValidator[T]) WithPrinterFunc(fn func(w io.Writer, a []T)) *InValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *InError[T]) {
		fn(w, e.ValidValues)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *InValidator[T]) Validate(v any) error {
	s := v.(T)
	if !slices.Contains(r.a, s) {
		return &InError[T]{
			Value:       s,
			ValidValues: r.a,
			rule:        r,
		}
	}
	return nil
}

// InError reports an error is caused in InValidator.
type InError[T comparable] struct {
	Value       T
	ValidValues []T
	rule        *InValidator[T]
}

// Error implements the error interface.
func (e InError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &inErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type inErrorPrinter[T comparable] struct{}

func (inErrorPrinter[T]) Print(w io.Writer, e *InError[T]) {
	fmt.Fprintf(w, "must be a valid value in %v", e.ValidValues)
}

// InErrorPrinter is the interface that wraps Print method.
type InErrorPrinter[T comparable] interface {
	Printer[InError[T]]
}

var _ typedValidator[
	*InValidator[string],
	InError[string],
	InErrorPrinter[string],
] = (*InValidator[string])(nil)
