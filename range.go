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
	p   MinViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MinValidator[T]) WithPrinter(p MinViolationPrinter[T]) *MinValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MinValidator[T]) WithPrinterFunc(fn func(w io.Writer, min T)) *MinValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MinViolationError[T]) {
		fn(w, e.Min)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MinValidator[T]) Validate(v any) error {
	n := v.(T)
	if n < r.min {
		return &MinViolationError[T]{
			Value: n,
			Min:   r.min,
			rule:  r,
		}
	}
	return nil
}

// MinViolationError reports an error is caused in MinValidator.
type MinViolationError[T ordered] struct {
	Value T
	Min   T
	rule  *MinValidator[T]
}

// Error implements the error interface.
func (e MinViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type minViolationPrinter[T ordered] struct{}

func (minViolationPrinter[T]) Print(w io.Writer, e *MinViolationError[T]) {
	fmt.Fprintf(w, "must be no less than %v", e.Min)
}

// MinViolationPrinter is the interface that wraps Print method.
type MinViolationPrinter[T ordered] interface {
	Printer[MinViolationError[T]]
}

var _ typedValidator[
	*MinValidator[int],
	MinViolationError[int],
	MinViolationPrinter[int],
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
	p   MaxViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *MaxValidator[T]) WithPrinter(p MaxViolationPrinter[T]) *MaxValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *MaxValidator[T]) WithPrinterFunc(fn func(w io.Writer, max T)) *MaxValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *MaxViolationError[T]) {
		fn(w, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *MaxValidator[T]) Validate(v any) error {
	n := v.(T)
	if n > r.max {
		return &MaxViolationError[T]{
			Value: n,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// MaxViolationError reports an error is caused in MaxValidator.
type MaxViolationError[T ordered] struct {
	Value T
	Max   T
	rule  *MaxValidator[T]
}

// Error implements the error interface.
func (e MaxViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type maxViolationPrinter[T ordered] struct{}

func (maxViolationPrinter[T]) Print(w io.Writer, e *MaxViolationError[T]) {
	fmt.Fprintf(w, "must be no greater than %v", e.Max)
}

// MaxViolationPrinter is the interface that wraps Print method.
type MaxViolationPrinter[T ordered] interface {
	Printer[MaxViolationError[T]]
}

var _ typedValidator[
	*MaxValidator[int],
	MaxViolationError[int],
	MaxViolationPrinter[int],
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
	p        InRangeViolationPrinter[T]
}

// WithPrinter returns shallow copy of r with its Printer changed to p.
func (r *InRangeValidator[T]) WithPrinter(p InRangeViolationPrinter[T]) *InRangeValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

// WithPrinterFunc returns shallow copy of r with its printer function changed to fn.
func (r *InRangeValidator[T]) WithPrinterFunc(fn func(w io.Writer, min, max T)) *InRangeValidator[T] {
	rr := *r
	rr.p = makePrinterFunc(func(w io.Writer, e *InRangeViolationError[T]) {
		fn(w, e.Min, e.Max)
	})
	return &rr
}

// Validate validates v. If v's type is not T, Validate panics.
func (r *InRangeValidator[T]) Validate(v any) error {
	n := v.(T)
	if n < r.min || n > r.max {
		return &InRangeViolationError[T]{
			Value: n,
			Min:   r.min,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

// InRangeViolationError reports an error is caused in InRangeValidator.
type InRangeViolationError[T ordered] struct {
	Value    T
	Min, Max T
	rule     *InRangeValidator[T]
}

// Error implements the error interface.
func (e InRangeViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &inRangeViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type inRangeViolationPrinter[T ordered] struct{}

func (inRangeViolationPrinter[T]) Print(w io.Writer, e *InRangeViolationError[T]) {
	fmt.Fprintf(w, "must be in range(%v ... %v)", e.Min, e.Max)
}

// InRangeViolationPrinter is the interface that wraps Print method.
type InRangeViolationPrinter[T ordered] interface {
	Printer[InRangeViolationError[T]]
}

var _ typedValidator[
	*InRangeValidator[int],
	InRangeViolationError[int],
	InRangeViolationPrinter[int],
] = (*InRangeValidator[int])(nil)
