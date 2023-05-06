package validator

import (
	"bytes"
	"errors"
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
		p = &structRulePrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type structRulePrinter[T any] struct{}

func (structRulePrinter[T]) Print(w io.Writer, e StructRuleViolationError[T]) {
	for _, err := range e.Errors {
		fmt.Fprintln(w, err)
	}
}

type StructRuleViolationPrinter[T any] interface {
	Printer[StructRuleViolationError[T]]
}

var (
	_ Validator                            = (*StructRuleValidator[struct{}])(nil)
	_ ViolationError                       = (*StructRuleViolationError[struct{}])(nil)
	_ StructRuleViolationPrinter[struct{}] = (*structRulePrinter[struct{}])(nil)
)

type StructRuleAdder interface {
	Add(field StructField, vs ...Validator)
}

type structRule[T any] struct {
	base   *T
	fields map[string]*structFieldRuleValidator
}

func (r *structRule[T]) Add(field StructField, vs ...Validator) {
	offset := field.offsetOf(r.base)
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

var _ StructRuleAdder = (*structRule[struct{}])(nil)

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
}

func (e StructFieldRuleViolationError[T]) Error() string {
	p := &structFieldRulePrinter[T]{}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type structFieldRulePrinter[T any] struct{}

func (structFieldRulePrinter[T]) Print(w io.Writer, e StructFieldRuleViolationError[T]) {
	errs, ok := e.Err.(interface{ Unwrap() []error })
	if !ok {
		// The structFieldRuleValidator.Validate always returns any errors
		// with errors.Join.
		panic("unexpected")
	}
	for i, err := range errs.Unwrap() {
		if i > 0 {
			w.Write([]byte("\n"))
		}
		fmt.Fprintf(w, "%s: %v", e.Name, err)
	}
}

type StructFieldRuleViolationPrinter[T any] interface {
	Printer[StructFieldRuleViolationError[T]]
}

var (
	_ Validator                                 = (*structFieldRuleValidator)(nil)
	_ ViolationError                            = (*StructFieldRuleViolationError[struct{}])(nil)
	_ StructFieldRuleViolationPrinter[struct{}] = (*structFieldRulePrinter[struct{}])(nil)
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
