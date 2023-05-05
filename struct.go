package validator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func Struct[T any](build func(s StructRuleAdder, v *T)) Validator {
	var (
		s structRuleValidator[T]
		v T
	)
	s.base = &v
	build(&s, &v)
	return &s
}

type StructRuleAdder interface {
	Add(field StructField, vs ...Validator)
}

type structRuleValidator[T any] struct {
	base  *T
	rules map[string]*structFieldRuleValidator
}

func (s *structRuleValidator[T]) Add(field StructField, vs ...Validator) {
	offset := field.offsetOf(s.base)
	f := lookupStructField(s.base, offset)
	if s.rules == nil {
		s.rules = make(map[string]*structFieldRuleValidator)
	}
	s.rules[field.Name()] = &structFieldRuleValidator{
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

func (s *structRuleValidator[T]) Validate(v any) error {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	base := p.Interface()
	errs := make(map[string]error)
	for name, rule := range s.rules {
		if err := rule.Validate(base); err != nil {
			errs[name] = err
		}
	}
	if len(errs) > 0 {
		return &StructRuleViolationError[T]{
			Value:  s.base,
			Errors: errs,
		}
	}
	return nil
}

type StructRuleViolationError[T any] struct {
	Value  *T
	Errors map[string]error
}

func (e StructRuleViolationError[T]) Error() string {
	var w bytes.Buffer
	for _, err := range e.Errors {
		fmt.Fprintln(&w, err)
	}
	return w.String()
}

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
		return r.field.createError(p, errors.Join(errs...))
	}
	return nil
}

type StructFieldRuleViolationError[T any] struct {
	Name  string
	Value T
	Err   error
	rule  *structRuleValidator[T]
}

func (e StructFieldRuleViolationError[T]) Error() string {
	p := &structFieldRulePrinter[T]{}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type structFieldRulePrinter[T any] struct{}

func (structFieldRulePrinter[T]) Print(w io.Writer, e StructFieldRuleViolationError[T]) {
	errs, ok := unwrapErrors(e.Err)
	if !ok {
		panic("unexpected")
	}
	for i, err := range errs {
		if i > 0 {
			w.Write([]byte("\n"))
		}
		fmt.Fprintf(w, "the field '%s' %v", e.Name, err)
	}
}

type StructRuleViolationPrinter[T any] interface {
	Printer[StructRuleViolationError[T]]
}

var (
	_ StructRuleAdder = (*structRuleValidator[struct{}])(nil)
	_ Validator       = (*structRuleValidator[struct{}])(nil)
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

func (f *structField[T]) offsetOf(base any) uintptr {
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
	offsetOf(base any) uintptr
	valueOf(base any, index []int) any
	createError(v any, err error) ViolationError
}

var _ StructField = (*structField[string])(nil)
