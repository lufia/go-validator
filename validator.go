// Package validator provides utilities for validating any types.
package validator

import (
	"errors"
	"io"
)

// Validator is the interface that wraps the basic Validate method.
type Validator interface {
	Validate(v any) error
}

type ViolationError interface {
	error
}

type Printer[E ViolationError] interface {
	Print(w io.Writer, e *E)
}

type PrinterFunc[E ViolationError] func(w io.Writer, e *E)

func printerFunc[E ViolationError](fn func(w io.Writer, e *E)) PrinterFunc[E] {
	return PrinterFunc[E](fn)
}

func (p PrinterFunc[E]) Print(w io.Writer, e *E) {
	p(w, e)
}

var _ Printer[RequiredViolationError[string]] = (PrinterFunc[RequiredViolationError[string]])(nil)

func Join(vs ...Validator) Validator {
	var a []Validator
	for _, v := range vs {
		if p, ok := v.(*joinValidator); ok {
			a = append(a, p.vs...)
		} else {
			a = append(a, v)
		}
	}
	return &joinValidator{vs: a}
}

type joinValidator struct {
	vs []Validator
}

func (r *joinValidator) Validate(v any) error {
	var errs []error
	for _, p := range r.vs {
		if err := p.Validate(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

var _ Validator = (*joinValidator)(nil)

type typedValidator[V Validator, E ViolationError, P Printer[E]] interface {
	Validator
	WithPrinter(p P) V
}
