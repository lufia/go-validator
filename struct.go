package validator

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func Struct[T any](build func(s StructRuleAdder, v *T)) *StructRuleValidator[T] {
	var (
		s StructRuleValidator[T]
		v T
	)
	rule := structRule[T]{
		base:   &v,
		fields: make(map[string]*structFieldRuleValidator),
	}
	build(&rule, &v)
	s.rule = &rule
	return &s
}

type StructRuleValidator[T any] struct {
	rule *structRule[T]
	p    StructRuleViolationPrinter[T]
}

func (r *StructRuleValidator[T]) WithPrinter(p StructRuleViolationPrinter[T]) *StructRuleValidator[T] {
	rr := *r
	rr.p = p
	return &rr
}

func (r *StructRuleValidator[T]) Validate(v any) error {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	base := p.Interface()
	errs := make(map[string]error)
	for name, rule := range r.rule.fields {
		if err := rule.Validate(base); err != nil {
			errs[name] = err
		}
	}
	if len(errs) > 0 {
		return &StructRuleViolationError[T]{
			Value:  r.rule.base,
			Errors: errs,
			rule:   r,
		}
	}
	return nil
}

type StructRuleViolationError[T any] struct {
	Value  *T
	Errors map[string]error
	rule   *StructRuleValidator[T]
}

func (e StructRuleViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &structRuleViolationPrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

func (e StructRuleViolationError[T]) Unwrap() []error {
	if len(e.Errors) == 0 {
		return nil
	}
	errs := make([]error, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err)
	}
	return errs
}

type structRuleViolationPrinter[T any] struct{}

func (structRuleViolationPrinter[T]) Print(w io.Writer, e *StructRuleViolationError[T]) {
	for _, err := range e.Errors {
		fmt.Fprintln(w, err)
	}
}

type StructRuleViolationPrinter[T any] interface {
	Printer[StructRuleViolationError[T]]
}

var _ typedValidator[
	*StructRuleValidator[any],
	StructRuleViolationError[any],
	StructRuleViolationPrinter[any],
] = (*StructRuleValidator[any])(nil)

type StructRuleAdder interface {
	Add(field StructField, vs ...Validator)
}

type structRule[T any] struct {
	base   *T
	fields map[string]*structFieldRuleValidator
}

func (r *structRule[T]) Add(field StructField, vs ...Validator) {
	offset := field.offsetFrom(r.base)
	f := lookupStructField(r.base, offset)
	r.fields[field.Name()] = &structFieldRuleValidator{
		field: field,
		vs:    vs,
		index: f.Index,
	}
}

func lookupStructField(p any, offset uintptr) reflect.StructField {
	v := reflect.ValueOf(p)
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, f := range fields {
		if f.Offset == offset {
			return f
		}
	}
	panic("the pointer refers out of the struct")
}

var _ StructRuleAdder = (*structRule[any])(nil)

type structFieldRuleValidator struct {
	field StructField
	vs    []Validator
	index []int
}

func (r *structFieldRuleValidator) Validate(base any) error {
	p := r.field.valueOf(base, r.index)

	var errs []error
	for _, v := range r.vs {
		if err := v.Validate(p); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return r.field.createError(p, joinErrors(errs...))
	}
	return nil
}

type StructFieldRuleViolationError[T any] struct {
	Name  string
	Value T
	Err   error
}

func (e StructFieldRuleViolationError[T]) Error() string {
	p := &structFieldRuleViolationPrinter[T]{}
	var w bytes.Buffer
	p.Print(&w, &e)
	return w.String()
}

type structFieldRuleViolationPrinter[T any] struct{}

func (structFieldRuleViolationPrinter[T]) Print(w io.Writer, e *StructFieldRuleViolationError[T]) {
	for i, err := range flattenErrors(e.Err) {
		if i > 0 {
			w.Write([]byte("\n"))
		}
		fmt.Fprintf(w, "%s: %v", e.Name, err)
	}
}

func flattenErrors(err error) []error {
	var errs []error
	e, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return []error{err}
	}
	for _, err := range e.Unwrap() {
		errs = append(errs, flattenErrors(err)...)
	}
	return errs
}

type StructFieldRuleViolationPrinter[T any] interface {
	Printer[StructFieldRuleViolationError[T]]
}

// structFieldRuleValidator is not satisfy typedValidator.
// Because it does not implement WithPrinter(p Printer[E]) method.
var (
	_ Validator                            = (*structFieldRuleValidator)(nil)
	_ ViolationError                       = (*StructFieldRuleViolationError[any])(nil)
	_ StructFieldRuleViolationPrinter[any] = (*structFieldRuleViolationPrinter[any])(nil)
)

func Field[T any](p *T, name string, opts ...any) StructField {
	return &structField[T]{
		name: name,
		p:    p,
	}
}

type structField[T any] struct {
	name string
	p    *T
}

func (f *structField[T]) Name() string {
	return f.name
}

func (f *structField[T]) offsetFrom(base any) uintptr {
	bp := reflect.ValueOf(base).Pointer()
	pp := reflect.ValueOf(f.p).Pointer()
	return pp - bp
}

func (f *structField[T]) valueOf(base any, index []int) any {
	p := reflect.ValueOf(base)
	return p.FieldByIndex(index).Interface()
}

func (f *structField[T]) createError(v any, err error) ViolationError {
	return &StructFieldRuleViolationError[T]{
		Name:  f.name,
		Value: v.(T),
		Err:   err,
	}
}

type StructField interface {
	Name() string
	offsetFrom(base any) uintptr
	valueOf(base any, index []int) any
	createError(v any, err error) ViolationError
}

var _ StructField = (*structField[string])(nil)
