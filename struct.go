package requiring

import (
	"errors"
	"reflect"
)

func Struct[T any](build func(s StructRuleAdder, v *T)) Validator {
	var (
		s structRuleSet
		v T
	)
	s.base = &v
	build(&s, &v)
	return &s
}

type StructRuleAdder interface {
	Add(p any, name string, vs ...Validator)
}

type structRuleSet struct {
	name  string
	base  any
	rules map[string]*structRule
}

type structRule struct {
	vs []Validator

	offset uintptr // offset within struct, in bytes
	index  []int   // index sequence for reflect.Type.FieldByIndex
}

func (s *structRuleSet) Add(p any, name string, vs ...Validator) {
	off := offsetOf(s.base, p)
	f := lookupStructField(s.base, off)
	if s.rules == nil {
		s.rules = make(map[string]*structRule)
	}
	setNameAll(name, vs)
	s.rules[name] = &structRule{
		vs:     vs,
		offset: f.Offset,
		index:  f.Index,
	}
}

func setNameAll(name string, vs []Validator) {
	for _, v := range vs {
		v.SetName(name)
	}
}

func offsetOf(base, p any) uintptr {
	bp := reflect.ValueOf(base).Pointer()
	pp := reflect.ValueOf(p).Pointer()
	return pp - bp
}

func lookupStructField(p any, off uintptr) reflect.StructField {
	v := reflect.ValueOf(p)
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, f := range fields {
		if f.Offset == off {
			return f
		}
	}
	panic("xxx")
}

func (s *structRuleSet) SetName(name string) {
	s.name = name
}

func (s *structRuleSet) Validate(v any) error {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	var errs []error
	for _, rule := range s.rules {
		f := p.FieldByIndex(rule.index)
		if err := validateAll(rule.vs, f.Interface()); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func validateAll(vs []Validator, v any) error {
	var errs []error
	for _, p := range vs {
		if err := p.Validate(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

var (
	_ StructRuleAdder = (*structRuleSet)(nil)
	_ Validator       = (*structRuleSet)(nil)
)
