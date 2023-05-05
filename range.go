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

func Min[T ordered](n T, opts ...any) Validator {
	var r minValidator[T]
	r.min = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case MinViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type minValidator[T ordered] struct {
	min T
	p   MinViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *minValidator[T]) Validate(v any) error {
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
	rule  *minValidator[T]
}

func (e MinViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type minPrinter[T ordered] struct{}

func (minPrinter[T]) Print(w io.Writer, e MinViolationError[T]) {
	fmt.Fprintf(w, "must be no less than %v", e.Min)
}

type MinViolationPrinter[T ordered] interface {
	Printer[MinViolationError[T]]
}

var (
	_ Validator                = (*minValidator[int])(nil)
	_ ViolationError           = (*MinViolationError[int])(nil)
	_ MinViolationPrinter[int] = (*minPrinter[int])(nil)
)

func Max[T ordered](n T, opts ...any) Validator {
	var r maxValidator[T]
	r.max = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case MaxViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type maxValidator[T ordered] struct {
	max T
	p   MaxViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *maxValidator[T]) Validate(v any) error {
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
	rule  *maxValidator[T]
}

func (e MaxViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type maxPrinter[T ordered] struct{}

func (maxPrinter[T]) Print(w io.Writer, e MaxViolationError[T]) {
	fmt.Fprintf(w, "must be no greater than %v", e.Max)
}

type MaxViolationPrinter[T ordered] interface {
	Printer[MaxViolationError[T]]
}

var (
	_ Validator                = (*maxValidator[int])(nil)
	_ ViolationError           = (*MaxViolationError[int])(nil)
	_ MaxViolationPrinter[int] = (*maxPrinter[int])(nil)
)

func InRange[T ordered](min, max T, opts ...any) Validator {
	var r inRangeValidator[T]
	r.min = min
	r.max = max
	for _, opt := range opts {
		switch v := opt.(type) {
		case InRangeViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type inRangeValidator[T ordered] struct {
	min, max T
	p        InRangeViolationPrinter[T]
	pp       InvalidTypePrinter
}

func (r *inRangeValidator[T]) Validate(v any) error {
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
	rule     *inRangeValidator[T]
}

func (e InRangeViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &inRangePrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type inRangePrinter[T ordered] struct{}

func (inRangePrinter[T]) Print(w io.Writer, e InRangeViolationError[T]) {
	fmt.Fprintf(w, "must be in range(%v ... %v)", e.Min, e.Max)
}

type InRangeViolationPrinter[T ordered] interface {
	Printer[InRangeViolationError[T]]
}

var (
	_ Validator                    = (*inRangeValidator[int])(nil)
	_ ViolationError               = (*InRangeViolationError[int])(nil)
	_ InRangeViolationPrinter[int] = (*inRangePrinter[int])(nil)
)
