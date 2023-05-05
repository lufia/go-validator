package validator

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func MinLength[T ~string](n int, opts ...any) Validator {
	var r minLengthValidator[T]
	r.min = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case MinLengthViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type minLengthValidator[T ~string] struct {
	min int
	p   MinLengthViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *minLengthValidator[T]) Validate(v any) error {
	s, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(s),
			p:     r.pp,
		}
	}
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

type MinLengthViolationError[T ~string] struct {
	Value T
	Min   int
	rule  *minLengthValidator[T]
}

func (e MinLengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &minLengthPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type minLengthPrinter[T ~string] struct{}

func (minLengthPrinter[T]) Print(w io.Writer, e MinLengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be no less than %v", e.Min)
}

type MinLengthViolationPrinter[T ~string] interface {
	Printer[MinLengthViolationError[T]]
}

var (
	_ Validator                         = (*minLengthValidator[string])(nil)
	_ ViolationError                    = (*MinLengthViolationError[string])(nil)
	_ MinLengthViolationPrinter[string] = (*minLengthPrinter[string])(nil)
)

func MaxLength[T ~string](n int, opts ...any) Validator {
	var r maxLengthValidator[T]
	r.max = n
	for _, opt := range opts {
		switch v := opt.(type) {
		case MaxLengthViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type maxLengthValidator[T ~string] struct {
	max int
	p   MaxLengthViolationPrinter[T]
	pp  InvalidTypePrinter
}

func (r *maxLengthValidator[T]) Validate(v any) error {
	s, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(s),
			p:     r.pp,
		}
	}
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

type MaxLengthViolationError[T ~string] struct {
	Value T
	Max   int
	rule  *maxLengthValidator[T]
}

func (e MaxLengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &maxLengthPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type maxLengthPrinter[T ~string] struct{}

func (maxLengthPrinter[T]) Print(w io.Writer, e MaxLengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be no greater than %v", e.Max)
}

type MaxLengthViolationPrinter[T ~string] interface {
	Printer[MaxLengthViolationError[T]]
}

var (
	_ Validator                         = (*maxLengthValidator[string])(nil)
	_ ViolationError                    = (*MaxLengthViolationError[string])(nil)
	_ MaxLengthViolationPrinter[string] = (*maxLengthPrinter[string])(nil)
)

func Length[T ~string](min, max int, opts ...any) Validator {
	var r lengthValidator[T]
	r.min = min
	r.max = max
	for _, opt := range opts {
		switch v := opt.(type) {
		case LengthViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type lengthValidator[T ~string] struct {
	min, max int
	p        LengthViolationPrinter[T]
	pp       InvalidTypePrinter
}

func (r *lengthValidator[T]) Validate(v any) error {
	s, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Value: v,
			Type:  reflect.TypeOf(s),
			p:     r.pp,
		}
	}
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

type LengthViolationError[T ~string] struct {
	Value    T
	Min, Max int
	rule     *lengthValidator[T]
}

func (e LengthViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &lengthPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type lengthPrinter[T ~string] struct{}

func (lengthPrinter[T]) Print(w io.Writer, e LengthViolationError[T]) {
	fmt.Fprintf(w, "the length must be in range(%v ... %v)", e.Min, e.Max)
}

type LengthViolationPrinter[T ~string] interface {
	Printer[LengthViolationError[T]]
}

var (
	_ Validator                      = (*lengthValidator[string])(nil)
	_ ViolationError                 = (*LengthViolationError[string])(nil)
	_ LengthViolationPrinter[string] = (*lengthPrinter[string])(nil)
)
