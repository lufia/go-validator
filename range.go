package validator

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/exp/constraints"
)

type ordered interface {
	constraints.Ordered
}

func Min[T ordered](n T, opts ...any) *MinValidator[T] {
	var r MinValidator[T]
	r.min = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type MinValidator[T ordered] struct {
	min T
	p   MinViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *MinValidator[T]) WithPrinter(p MinViolationPrinter[T]) *MinValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

func (r *MinValidator[T]) WithPrinterFunc(fn func(w io.Writer, min T)) *MinValidator[T] {
	rr := *r
	rr.p = printerFunc(func(w io.Writer, e *MinViolationError[T]) {
		fn(w, e.Min)
	})
	return &rr
}

func (r *MinValidator[T]) Validate(v any) error {
	n, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(n),
			p:     r.pp,
		}
	}
	if n < r.min {
		return &MinViolationError[T]{
			Value: n,
			Min:   r.min,
			rule:  r,
		}
	}
	return nil
}

type MinViolationError[T ordered] struct {
	Value T
	Min   T
	rule  *MinValidator[T]
}

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

type MinViolationPrinter[T ordered] interface {
	Printer[MinViolationError[T]]
}

var _ typedValidator[
	*MinValidator[int],
	MinViolationError[int],
	MinViolationPrinter[int],
] = (*MinValidator[int])(nil)

func Max[T ordered](n T, opts ...any) *MaxValidator[T] {
	var r MaxValidator[T]
	r.max = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type MaxValidator[T ordered] struct {
	max T
	p   MaxViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *MaxValidator[T]) WithPrinter(p MaxViolationPrinter[T]) *MaxValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

func (r *MaxValidator[T]) WithPrinterFunc(fn func(w io.Writer, max T)) *MaxValidator[T] {
	rr := *r
	rr.p = printerFunc(func(w io.Writer, e *MaxViolationError[T]) {
		fn(w, e.Max)
	})
	return &rr
}

func (r *MaxValidator[T]) Validate(v any) error {
	n, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(n),
			p:     r.pp,
		}
	}
	if n > r.max {
		return &MaxViolationError[T]{
			Value: n,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

type MaxViolationError[T ordered] struct {
	Value T
	Max   T
	rule  *MaxValidator[T]
}

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

type MaxViolationPrinter[T ordered] interface {
	Printer[MaxViolationError[T]]
}

var _ typedValidator[
	*MaxValidator[int],
	MaxViolationError[int],
	MaxViolationPrinter[int],
] = (*MaxValidator[int])(nil)

func InRange[T ordered](min, max T, opts ...any) *InRangeValidator[T] {
	var r InRangeValidator[T]
	r.min = min
	r.max = max
	for _, opt := range opts {
		switch v := opt.(type) {
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type InRangeValidator[T ordered] struct {
	min, max T
	p        InRangeViolationPrinter[T]
	pp       InvalidTypePrinter
}

func (r *InRangeValidator[T]) WithPrinter(p InRangeViolationPrinter[T]) *InRangeValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

func (r *InRangeValidator[T]) WithPrinterFunc(fn func(w io.Writer, min, max T)) *InRangeValidator[T] {
	rr := *r
	rr.p = printerFunc(func(w io.Writer, e *InRangeViolationError[T]) {
		fn(w, e.Min, e.Max)
	})
	return &rr
}

func (r *InRangeValidator[T]) Validate(v any) error {
	n, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(n),
			p:     r.pp,
		}
	}
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

type InRangeViolationError[T ordered] struct {
	Value    T
	Min, Max T
	rule     *InRangeValidator[T]
}

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

type InRangeViolationPrinter[T ordered] interface {
	Printer[InRangeViolationError[T]]
}

var _ typedValidator[
	*InRangeValidator[int],
	InRangeViolationError[int],
	InRangeViolationPrinter[int],
] = (*InRangeValidator[int])(nil)
