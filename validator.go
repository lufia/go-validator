// Package requiring provides utilities for validating any types.
package requiring

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

type ViolationError interface {
	error
}

type Printer[E ViolationError] interface {
	Print(w io.Writer, e E)
}

type Validator interface {
	Validate(v any) error
}

type rule struct {
	Name       string
	Validators []Validator
	Offset     uintptr // offset within struct, in bytes
	Index      []int   // index sequence for reflect.Type.FieldByIndex
}

func (r *rule) Validate(v any) error {
	var errs []error
	for _, p := range r.Validators {
		if err := p.Validate(v); err != nil {
			errs = append(errs, fmt.Errorf("'%s' %w", r.Name, err))
		}
	}
	return errors.Join(errs...)
}

type RuleSet struct {
	base  any
	rules map[string]*rule
}

func (s *RuleSet) Add(p any, name string, vs ...Validator) {
	off := s.offsetOf(p)
	f := lookupStructField(s.base, off)
	if s.rules == nil {
		s.rules = make(map[string]*rule)
	}
	s.rules[name] = &rule{
		Name:       name,
		Validators: vs,
		Offset:     f.Offset,
		Index:      f.Index,
	}
}

func (s *RuleSet) offsetOf(p any) uintptr {
	bp := reflect.ValueOf(s.base).Pointer()
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

func (s *RuleSet) Validate(v any) error {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	var errs []error
	for _, rule := range s.rules {
		f := p.FieldByIndex(rule.Index)
		if err := rule.Validate(f.Interface()); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func Struct[T any](build func(s *RuleSet, v *T)) Validator {
	var (
		s RuleSet
		v T
	)
	s.base = &v
	build(&s, &v)
	return &s
}
