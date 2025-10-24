package validator

import (
	"context"

	"golang.org/x/text/message"
)

// Pointer returns the validator to verify a pointer to something.
func Pointer[T any](vs ...Validator[T]) Validator[*T] {
	return &pointerValidator[*T, T]{
		vs: vs,
	}
}

// pointerValidator represents the validator to check pointer value.
type pointerValidator[P ~*T, T any] struct {
	vs []Validator[T]
}

// WithFormat returns shallow copy of r with its error format changed to key.
//
// TODO(lufia): currently key is always ignored.
func (r *pointerValidator[P, T]) WithFormat(key message.Reference, a ...Arg) Validator[P] {
	rr := *r
	return &rr
}

// Validate validates v.
func (r *pointerValidator[P, T]) Validate(ctx context.Context, p P) error {
	if p == nil {
		return nil
	}
	var errs []error
	for _, rule := range r.vs {
		if err := rule.Validate(ctx, *p); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return joinErrors(errs...)
	}
	return nil
}
