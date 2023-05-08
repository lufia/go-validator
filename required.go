package validator

import (
	"bytes"
	"fmt"
	"io"
)

// Required returns the validator to verify the value is not zero value.
func Required[T comparable]() *RequiredValidator[T] {
	return &RequiredValidator[T]{}
}

type RequiredValidator[T comparable] struct {
	p RequiredViolationPrinter[T]
}

func (r *RequiredValidator[T]) WithPrinter(p RequiredViolationPrinter[T]) *RequiredValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

func (r *RequiredValidator[T]) WithPrinterFunc(fn func(w io.Writer)) *RequiredValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, _ *RequiredViolationError[T]) {
		fn(w)
	})
	return &rr
}

func (r *RequiredValidator[T]) Validate(v any) error {
	s := v.(T)
	var v0 T
	if s == v0 {
		return &RequiredViolationError[T]{
			Value: s,
			rule:  r,
		}
	}
	return nil
}

type RequiredViolationError[T comparable] struct {
	Value T
	rule  *RequiredValidator[T]
}

func (e RequiredViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &requiredPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type requiredPrinter[T comparable] struct{}

func (requiredPrinter[T]) Print(w io.Writer, e *RequiredViolationError[T]) {
	fmt.Fprintf(w, "cannot be the zero value")
}

type RequiredViolationPrinter[T comparable] interface {
	Printer[RequiredViolationError[T]]
}

var _ typedValidator[
	*RequiredValidator[string],
	RequiredViolationError[string],
	RequiredViolationPrinter[string],
] = (*RequiredValidator[string])(nil)
