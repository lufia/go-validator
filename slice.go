package validator

import (
	"context"

	"golang.org/x/text/message"
)

type slice[T any] interface {
	[]T
}

func Slice[T any](vs ...Validator[T]) Validator[[]T] {
	return &sliceValidator[[]T, T]{
		vs: vs,
	}
}

// sliceValidator represents the validator to check slice elements.
type sliceValidator[S slice[T], T any] struct {
	vs []Validator[T]
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
//
// TODO(lufia): currently key is always ignored.
func (r *sliceValidator[S, T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[S] {
	rr := *r
	return &rr
}

// Validate validates v.
func (r *sliceValidator[S, T]) Validate(ctx context.Context, v S) error {
	var m OrderedMap[int, error]
	for i, elem := range v {
		var errs []error
		for _, rule := range r.vs {
			if err := rule.Validate(ctx, elem); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			m.set(i, joinErrors(errs...))
		}
	}
	if m.Len() > 0 {
		return &SliceError[S, T]{
			Value:  v,
			Errors: &m,
		}
	}
	return nil
}

// SliceError reports an error is caused in Slice validator.
type SliceError[S slice[T], T any] struct {
	Value  S
	Errors *OrderedMap[int, error]
}

// Error implements the error interface.
func (e SliceError[S, T]) Error() string {
	return joinErrors(e.Unwrap()...).Error()
}

// Unwrap returns each errors of err.
func (e SliceError[S, T]) Unwrap() []error {
	n := e.Errors.Len()
	if n == 0 {
		return nil
	}
	errs := make([]error, 0, n)
	for _, key := range e.Errors.Keys() {
		err, _ := e.Errors.Get(key)
		errs = append(errs, err)
	}
	return errs
}

var (
	_ Validator[[]any] = (*sliceValidator[[]any, any])(nil)
	_ Error            = (*SliceError[[]any, any])(nil)
)
