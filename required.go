package requiring

import (
	"bytes"
	"fmt"
	"io"
)

func Required[T comparable](opts ...any) Validator {
	var r requiredValidator[T]
	for _, opt := range opts {
		switch v := opt.(type) {
		case RequiredViolationPrinter[T]:
			r.p = v
		}
	}
	return &r
}

type requiredValidator[T comparable] struct {
	name string
	p    RequiredViolationPrinter[T]
}

func (r *requiredValidator[T]) SetName(name string) {
	r.name = name
}

func (r *requiredValidator[T]) Validate(v any) error {
	s := v.(T)
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
		fmt.Fprintf(w, "'%s' is ", e.Name)
	}
	fmt.Fprintf(w, "required")
}

type RequiredViolationPrinter[T comparable] interface {
	Printer[RequiredViolationError[T]]
}

var (
	_ Validator                        = (*requiredValidator[string])(nil)
	_ ViolationError                   = (*RequiredViolationError[string])(nil)
	_ RequiredViolationPrinter[string] = (*requiredPrinter[string])(nil)
)
