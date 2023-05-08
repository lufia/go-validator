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
	p   MinLengthViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MinLengthValidator[T]) WithPrinter(p MinLengthViolationPrinter[T]) *MinLengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MinLengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, min int)) *MinLengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MinLengthViolationError[T]) {
		fn(w, e.Min)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MinLengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) < r.min {
		return &MinLengthViolationError[T]{
			Value: s,
			Min:   r.min,
			rule:  r,
		}
	}
	return nil
}

// MinLengthViolationError reports an error is caused in MinLengthValidator.
type MinLengthViolationError[T ~string] struct {
	Value T
	Min   int
	rule  *MinLengthValidator[T]
}

// Error implements the error interface.
func (e MinLengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minLengthViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type minLengthViolationPrinter[T ~string] struct{}

func (minLengthViolationPrinter[T]) Print(w io.Writer, e *MinLengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be no less than %v", e.Min)
}

// MinLengthViolationPrinter is the interface that wraps Print method.
type MinLengthViolationPrinter[T ~string] interface {
	Printer[MinLengthViolationError[T]]
}

var _ typedValidator[
	*MinLengthValidator[string],
	MinLengthViolationError[string],
	MinLengthViolationPrinter[string],
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
	p   MaxLengthViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MaxLengthValidator[T]) WithPrinter(p MaxLengthViolationPrinter[T]) *MaxLengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MaxLengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, max int)) *MaxLengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MaxLengthViolationError[T]) {
		fn(w, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MaxLengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) > r.max {
		return &MaxLengthViolationError[T]{
			Value: s,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// MaxLengthViolationError reports an error is caused in MaxLengthValidator.
type MaxLengthViolationError[T ~string] struct {
	Value T
	Max   int
	rule  *MaxLengthValidator[T]
}

// Error implements the error interface.
func (e MaxLengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxLengthViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type maxLengthViolationPrinter[T ~string] struct{}

func (maxLengthViolationPrinter[T]) Print(w io.Writer, e *MaxLengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be no greater than %v", e.Max)
}

// MaxLengthViolationPrinter is the interface that wraps Print method.
type MaxLengthViolationPrinter[T ~string] interface {
	Printer[MaxLengthViolationError[T]]
}

var _ typedValidator[
	*MaxLengthValidator[string],
	MaxLengthViolationError[string],
	MaxLengthViolationPrinter[string],
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
	p        LengthViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *LengthValidator[T]) WithPrinter(p LengthViolationPrinter[T]) *LengthValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *LengthValidator[T]) WithPrinterFunc(fn func(w io.Writer, min, max int)) *LengthValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *LengthViolationError[T]) {
		fn(w, e.Min, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *LengthValidator[T]) Validate(v any) error {
	s := v.(T)
	a := []rune(s)
	if len(a) < r.min || len(a) > r.max {
		return &LengthViolationError[T]{
			Value: s,
			Min:   r.min,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// LengthViolationError reports an error is caused in LengthValidator.
type LengthViolationError[T ~string] struct {
	Value    T
	Min, Max int
	rule     *LengthValidator[T]
}

// Error implements the error interface.
func (e LengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &lengthViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type lengthViolationPrinter[T ~string] struct{}

func (lengthViolationPrinter[T]) Print(w io.Writer, e *LengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be in range(%v ... %v)", e.Min, e.Max)
}

// LengthViolationPrinter is the interface that wraps Print method.
type LengthViolationPrinter[T ~string] interface {
	Printer[LengthViolationError[T]]
}

var _ typedValidator[
	*LengthValidator[string],
	LengthViolationError[string],
	LengthViolationPrinter[string],
] = (*LengthValidator[string])(nil)
