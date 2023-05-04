package requiring

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// Required returns the validator to verify the value is not zero value.
// When opts contains the type RequiredViolationPrinter[T],
// it will be used to print the required violation error.
// Also, when opts contains the type InvalidTypePrinter,
// it will be used to print the invalid type error.
func Required[T comparable](opts ...any) Validator {
	var r requiredValidator[T]
	for _, opt := range opts {
		switch v := opt.(type) {
		case RequiredViolationPrinter[T]:
			r.p = v
		case InvalidTypePrinter:
			r.pp = v
		}
	}
	return &r
}

type requiredValidator[T comparable] struct {
	name string
	p    RequiredViolationPrinter[T]
	pp   InvalidTypePrinter
}

func (r *requiredValidator[T]) SetName(name string) {
	r.name = name
}

func (r *requiredValidator[T]) Validate(v any) error {
	s, ok := v.(T)
	if !ok {
		return &InvalidTypeError{
			Name:  r.name,
			Value: v,
			Type:  reflect.TypeOf(s),
			p:     r.pp,
		}
	}
	var v0 T
	if s == v0 {
		return &RequiredViolationError[T]{
			Name:  r.name,
			Value: s,
			rule:  r,
		}
	}
	return nil
}

type RequiredViolationError[T comparable] struct {
	Name  string
	Value T
	rule  *requiredValidator[T]
}

func (e RequiredViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &requiredPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type requiredPrinter[T comparable] struct{}

func (requiredPrinter[T]) Print(w io.Writer, e RequiredViolationError[T]) {
	if e.Name != "" {
		fmt.Fprintf(w, "the field '%s' ", e.Name)
	}
	fmt.Fprintf(w, "cannot be the zero value")
}

type RequiredViolationPrinter[T comparable] interface {
	Printer[RequiredViolationError[T]]
}

var (
	_ Validator                        = (*requiredValidator[string])(nil)
	_ ViolationError                   = (*RequiredViolationError[string])(nil)
	_ RequiredViolationPrinter[string] = (*requiredPrinter[string])(nil)
)
