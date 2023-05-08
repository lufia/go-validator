package validator

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

type ordered interface {
	constraints.Ordered
}

// Min returns the validator to verify the value is greater or equal than n.
func Min[T ordered](n T) *MinValidator[T] {
	var r MinValidator[T]
	r.min = n
	return &r
}

// MinLengthValidator represents the validator to check the value is greater or equal than T.
type MinValidator[T ordered] struct {
	min T
	p   MinErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MinValidator[T]) WithPrinter(p MinErrorPrinter[T]) *MinValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MinValidator[T]) WithPrinterFunc(fn func(w io.Writer, min T)) *MinValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MinError[T]) {
		fn(w, e.Min)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MinValidator[T]) Validate(v any) error {
	n := v.(T)
	if n < r.min {
		return &MinError[T]{
			Value: n,
			Min:   r.min,
			rule:  r,
		}
	}
	return nil
}

// MinError reports an error is caused in MinValidator.
type MinError[T ordered] struct {
	Value T
	Min   T
	rule  *MinValidator[T]
}

// Error implements the error interface.
func (e MinError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type minErrorPrinter[T ordered] struct{}

func (minErrorPrinter[T]) Print(w io.Writer, e *MinError[T]) {
	fmt.Fprintf(w, "must be no less than %v", e.Min)
}

// MinErrorPrinter is the interface that wraps Print method.
type MinErrorPrinter[T ordered] interface {
	Printer[MinError[T]]
}

var _ typedValidator[
	*MinValidator[int],
	MinError[int],
	MinErrorPrinter[int],
] = (*MinValidator[int])(nil)

// Max returns the validator to verify the value is less or equal than n.
func Max[T ordered](n T) *MaxValidator[T] {
	var r MaxValidator[T]
	r.max = n
	return &r
}

// MaxValidator represents the validator to check the value is less or equal than T.
type MaxValidator[T ordered] struct {
	max T
	p   MaxErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MaxValidator[T]) WithPrinter(p MaxErrorPrinter[T]) *MaxValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MaxValidator[T]) WithPrinterFunc(fn func(w io.Writer, max T)) *MaxValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MaxError[T]) {
		fn(w, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MaxValidator[T]) Validate(v any) error {
	n := v.(T)
	if n > r.max {
		return &MaxError[T]{
			Value: n,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// MaxError reports an error is caused in MaxValidator.
type MaxError[T ordered] struct {
	Value T
	Max   T
	rule  *MaxValidator[T]
}

// Error implements the error interface.
func (e MaxError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type maxErrorPrinter[T ordered] struct{}

func (maxErrorPrinter[T]) Print(w io.Writer, e *MaxError[T]) {
	fmt.Fprintf(w, "must be no greater than %v", e.Max)
}

// MaxErrorPrinter is the interface that wraps Print method.
type MaxErrorPrinter[T ordered] interface {
	Printer[MaxError[T]]
}

var _ typedValidator[
	*MaxValidator[int],
	MaxError[int],
	MaxErrorPrinter[int],
] = (*MaxValidator[int])(nil)

// InRange returns the validator to verify the value is within min and max.
func InRange[T ordered](min, max T) *InRangeValidator[T] {
	var r InRangeValidator[T]
	r.min = min
	r.max = max
	return &r
}

// InRangeValidator represents the validator to check the value is within T.
type InRangeValidator[T ordered] struct {
	min, max T
	p        InRangeErrorPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *InRangeValidator[T]) WithPrinter(p InRangeErrorPrinter[T]) *InRangeValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *InRangeValidator[T]) WithPrinterFunc(fn func(w io.Writer, min, max T)) *InRangeValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *InRangeError[T]) {
		fn(w, e.Min, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *InRangeValidator[T]) Validate(v any) error {
	n := v.(T)
	if n < r.min || n > r.max {
		return &InRangeError[T]{
			Value: n,
			Min:   r.min,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// InRangeError reports an error is caused in InRangeValidator.
type InRangeError[T ordered] struct {
	Value    T
	Min, Max T
	rule     *InRangeValidator[T]
}

// Error implements the error interface.
func (e InRangeError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &inRangeErrorPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type inRangeErrorPrinter[T ordered] struct{}

func (inRangeErrorPrinter[T]) Print(w io.Writer, e *InRangeError[T]) {
	fmt.Fprintf(w, "must be in range(%v ... %v)", e.Min, e.Max)
}

// InRangeErrorPrinter is the interface that wraps Print method.
type InRangeErrorPrinter[T ordered] interface {
	Printer[InRangeError[T]]
}

var _ typedValidator[
	*InRangeValidator[int],
	InRangeError[int],
	InRangeErrorPrinter[int],
] = (*InRangeValidator[int])(nil)
