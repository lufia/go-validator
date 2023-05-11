// Package validator provides utilities for validating any types.
package validator

import (
	"io"
)

// Validator is the interface that wraps the basic Validate method.
type Validator interface {
	Validate(v any) error
}

// Error is the interface that wraps Error method.
type Error interface {
	error
}

// Printer is the interface that wraps Print method.
type Printer[E Error] interface {
	Print(w io.Writer, e *E)
}

type printerFunc[E Error] func(w io.Writer, e *E)

func makePrinterFunc[E Error](fn func(w io.Writer, e *E)) printerFunc[E] {
	return printerFunc[E](fn)
}

func (p printerFunc[E]) Print(w io.Writer, e *E) {
	p(w, e)
}

var _ Printer[RequiredError[string]] = (printerFunc[RequiredError[string]])(nil)

// Join bundles vs to a validator.
func Join(vs ...Validator) Validator {
	var a []Validator
	for _, v := range vs {
		if p, ok := v.(*joinValidator); ok {
			a = append(a, p.vs...)
		} else {
			a = append(a, v)
		}
	}
	return &joinValidator{vs: a}
}

type joinValidator struct {
	vs []Validator
}

// Validate returns the all errors that v is validated with its each validator.
func (r *joinValidator) Validate(v any) error {
	var errs []error
	for _, p := range r.vs {
		if err := p.Validate(v); err != nil {
			errs = append(errs, err)
		}
	}
	return joinErrors(errs...)
}

var _ Validator = (*joinValidator)(nil)

type typedValidator[V Validator, E Error, P Printer[E]] interface {
	Validator
	WithPrinter(p P) V
}

// OrderedMap is a map that guarantee that the iteration order of entries
// will be the order in which they were set.
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

func (m *OrderedMap[K, V]) set(key K, v V) {
	if m.values == nil {
		m.values = make(map[K]V)
	}
	if _, ok := m.values[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.values[key] = v
}

// Len returns the length of m.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

// Keys returns keys of m. These keys preserves the order in which they were set.
func (m *OrderedMap[K, V]) Keys() []K {
	return m.keys
}

// Get returns a value associated to key.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.values[key]
	return v, ok
}
