package validator

import (
	"bytes"
	"fmt"
	"io"
)

// MinLength returns the validator to verify the length of the value is greater or equal than n.
func MinLength[T ~string](n int) *MinLengthValidator[T] {
	var r MinLengthValidator[T]
	r.min = n
	return &r
}

// MinLengthValidator represents the validator to check the length of the value is greater or equal than T.
type MinLengthValidator[T ~string] struct {
	min int
	p   MinLengthErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MinLengthValidator[T]) WithPrinter(p MinLengthErrorPrinter[T]) *MinLengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MinLengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, min int)) *MinLengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MinLengthError[T]) {
		fn(w, e.Min)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MinLengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) < r.min {
		return &MinLengthError[T]{
			Value: s,
			Min:   r.min,
			rule:  r,
		}
	}
	return nil
}

// MinLengthError reports an error is caused in MinLengthValidator.
type MinLengthError[T ~string] struct {
	Value T
	Min   int
	rule  *MinLengthValidator[T]
}

// Error implements the error interface.
func (e MinLengthError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minLengthErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type minLengthErrorPrinter[T ~string] struct{}

func (minLengthErrorPrinter[T]) Print(w io.Writer, e *MinLengthError[T]) {
	fmt.Fprintf(w, "the length must be no less than %v", e.Min)
}

// MinLengthErrorPrinter is the interface that wraps Print method.
type MinLengthErrorPrinter[T ~string] interface {
	Printer[MinLengthError[T]]
}

var _ typedValidator[
	*MinLengthValidator[string],
	MinLengthError[string],
	MinLengthErrorPrinter[string],
] = (*MinLengthValidator[string])(nil)

// MaxLength returns the validator to verify the length of the value is less or equal than n.
func MaxLength[T ~string](n int) *MaxLengthValidator[T] {
	var r MaxLengthValidator[T]
	r.max = n
	return &r
}

// MaxLengthValidator represents the validator to check the length of the value is less or equal than T.
type MaxLengthValidator[T ~string] struct {
	max int
	p   MaxLengthErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MaxLengthValidator[T]) WithPrinter(p MaxLengthErrorPrinter[T]) *MaxLengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MaxLengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, max int)) *MaxLengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MaxLengthError[T]) {
		fn(w, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MaxLengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) > r.max {
		return &MaxLengthError[T]{
			Value: s,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// MaxLengthError reports an error is caused in MaxLengthValidator.
type MaxLengthError[T ~string] struct {
	Value T
	Max   int
	rule  *MaxLengthValidator[T]
}

// Error implements the error interface.
func (e MaxLengthError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxLengthErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type maxLengthErrorPrinter[T ~string] struct{}

func (maxLengthErrorPrinter[T]) Print(w io.Writer, e *MaxLengthError[T]) {
	fmt.Fprintf(w, "the length must be no greater than %v", e.Max)
}

// MaxLengthErrorPrinter is the interface that wraps Print method.
type MaxLengthErrorPrinter[T ~string] interface {
	Printer[MaxLengthError[T]]
}

var _ typedValidator[
	*MaxLengthValidator[string],
	MaxLengthError[string],
	MaxLengthErrorPrinter[string],
] = (*MaxLengthValidator[string])(nil)

// Length returns the validator to verify the length of the value is within min and max.
func Length[T ~string](min, max int) *LengthValidator[T] {
	var r LengthValidator[T]
	r.min = min
	r.max = max
	return &r
}

// LengthValidator represents the validator to check the length of the value is within T.
type LengthValidator[T ~string] struct {
	min, max int
	p        LengthErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *LengthValidator[T]) WithPrinter(p LengthErrorPrinter[T]) *LengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *LengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, min, max int)) *LengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *LengthError[T]) {
		fn(w, e.Min, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *LengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) < r.min || len(a) > r.max {
		return &LengthError[T]{
			Value: s,
			Min:   r.min,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// LengthError reports an error is caused in LengthValidator.
type LengthError[T ~string] struct {
	Value    T
	Min, Max int
	rule     *LengthValidator[T]
}

// Error implements the error interface.
func (e LengthError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &lengthErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type lengthErrorPrinter[T ~string] struct{}

func (lengthErrorPrinter[T]) Print(w io.Writer, e *LengthError[T]) {
	fmt.Fprintf(w, "the length must be in range(%v ... %v)", e.Min, e.Max)
}

// LengthErrorPrinter is the interface that wraps Print method.
type LengthErrorPrinter[T ~string] interface {
	Printer[LengthError[T]]
}

var _ typedValidator[
	*LengthValidator[string],
	LengthError[string],
	LengthErrorPrinter[string],
] = (*LengthValidator[string])(nil)
