package validator

import (
	"context"
	"errors"
	"reflect"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const structFieldKey = "%s: %v"

func init() {
	DefaultCatalog.SetString(language.English, structFieldKey, structFieldKey)
	DefaultCatalog.SetString(language.Japanese, structFieldKey, "xx")
}

// Struct returns the validator to verify that the struct satisfies rules constrated with build.
func Struct[T any](build func(s StructFieldAdder, v *T)) Validator[T] {
	var (
		s structValidator[T]
		v T
	)
	rule := structRule[T]{
		base:   &v,
		fields: make(map[string]StructField),
	}
	build(&rule, &v)
	s.rule = &rule
	s.key = structFieldKey
	s.args = []Arg{ByName("name"), ByName("error")}
	return &s
}

// structValidator represents the validator to check the struct satisfies its rules.
type structValidator[T any] struct {
	rule *structRule[T]
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *structValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	r.key = key
	r.args = a
	return &rr
}

// Validate validates v.
func (r *structValidator[T]) Validate(ctx context.Context, v T) error {
	errs := make(map[string]error)
	for name, rule := range r.rule.fields {
		if err := rule.validateField(ctx, &v, r.key, r.args); err != nil {
			errs[name] = err
		}
	}
	if len(errs) > 0 {
		return &StructError[T]{
			Value:  v,
			Errors: errs,
		}
	}
	return nil
}

// StructError reports an error is caused in Struct validator.
type StructError[T any] struct {
	Value  T
	Errors map[string]error
}

// Error implements the error interface.
func (e StructError[T]) Error() string {
	return joinErrors(e.Unwrap()...).Error()
}

// Unwrap returns each errors of err.
func (e StructError[T]) Unwrap() []error {
	if len(e.Errors) == 0 {
		return nil
	}
	errs := make([]error, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err)
	}
	return errs
}

var (
	_ Validator[any] = (*structValidator[any])(nil)
	_ Error          = (*StructError[any])(nil)
)

// StructFieldAdder is the interface that wraps Add method.
type StructFieldAdder interface {
	Add(field StructField)
}

type structRule[T any] struct {
	base   *T
	fields map[string]StructField
}

// Add adds the rule.
func (r *structRule[T]) Add(field StructField) {
	offset := field.offsetFrom(r.base)
	f := lookupStructField(r.base, offset)
	field.setIndex(f.Index)
	r.fields[field.Name()] = field
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

var _ StructFieldAdder = (*structRule[any])(nil)

// Field returns the p's field of the struct T.
func Field[T any](p *T, name string, vs ...Validator[T]) StructField {
	return &structField[T]{
		name: name,
		p:    p,
		vs:   vs,
	}
}

type structField[T any] struct {
	name  string
	p     *T
	vs    []Validator[T]
	index []int
}

// Name returns field's name.
func (r *structField[T]) Name() string {
	return r.name
}

func (r *structField[T]) setIndex(index []int) {
	if r.index != nil {
		panic("the field is already added")
	}
	r.index = index
}

func (r *structField[T]) validateField(ctx context.Context, base any, key message.Reference, args []Arg) error {
	p := r.valueOf(base, r.index)
	v := p.(T)
	var errs []error
	for _, rule := range r.vs {
		if err := rule.Validate(ctx, v); err != nil {
			err = wrapErrors(err, func(err error) error {
				e := &structFieldError[T]{
					Name:  r.name,
					Value: v,
					Err:   err,
				}
				return errors.New(ctxPrint(ctx, e, key, args))
			})
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return joinErrors(errs...)
	}
	return nil
}

func wrapErrors(err error, fn func(err error) error) error {
	errs := flattenErrors(err)
	for i, err := range errs {
		errs[i] = fn(err)
	}
	return joinErrors(errs...)
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

func (f *structField[T]) offsetFrom(base any) uintptr {
	bp := reflect.ValueOf(base).Pointer()
	pp := reflect.ValueOf(f.p).Pointer()
	return pp - bp
}

func (f *structField[T]) valueOf(base any, index []int) any {
	p := reflect.ValueOf(base).Elem()
	return p.FieldByIndex(index).Interface()
}

// structFieldError reports an error is caused in Field validator.
type structFieldError[T any] struct {
	Name  string `arg:"name"`
	Value T      `arg:"value"`
	Err   error  `arg:"error"`
}

// StructField is the interface that is used by StructFieldAdder.
type StructField interface {
	Name() string
	setIndex(index []int)
	validateField(ctx context.Context, base any, key message.Reference, a []Arg) error
	offsetFrom(base any) uintptr
}

var _ StructField = (*structField[string])(nil)
