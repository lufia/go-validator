package requiring

import (
	"bytes"
	"fmt"
	"io"
)

func NotEmpty[T ~string](opts ...any) Validator {
	var r notEmptyValidator[T]
	for _, opt := range opts {
		switch v := opt.(type) {
		case NotEmptyViolationPrinter[T]:
			r.p = v
		}
	}
	return &r
}

type notEmptyValidator[T ~string] struct {
	name string
	p    NotEmptyViolationPrinter[T]
}

func (r *notEmptyValidator[T]) SetName(name string) {
	r.name = name
}

func (r *notEmptyValidator[T]) Validate(v any) error {
	s := v.(T)
	if s == "" {
		return &NotEmptyViolationError[T]{
			Name:  r.name,
			Value: s,
			rule:  r,
		}
	}
	return nil
}

type NotEmptyViolationError[T ~string] struct {
	Name  string
	Value T
	rule  *notEmptyValidator[T]
	msg   string
}

func (e NotEmptyViolationError[T]) Error() string {
	if e.msg != "" {
		return e.msg
	}
	p := e.rule.p
	if p == nil {
		p = &notEmptyPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	e.msg = w.String()
	return e.msg
}

type notEmptyPrinter[T ~string] struct{}

func (notEmptyPrinter[T]) Print(w io.Writer, e NotEmptyViolationError[T]) {
	if e.Name != "" {
		fmt.Fprintf(w, "'%s' is ", e.Name)
	}
	fmt.Fprintf(w, "required")
}

type NotEmptyViolationPrinter[T ~string] interface {
	Printer[NotEmptyViolationError[T]]
}

var (
	_ Validator                        = (*notEmptyValidator[string])(nil)
	_ ViolationError                   = (*NotEmptyViolationError[string])(nil)
	_ NotEmptyViolationPrinter[string] = (*notEmptyPrinter[string])(nil)
)
