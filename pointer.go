package validator

import (
	"context"
)

// Pointer returns the validator to verify the pointer.
func Pointer[T any](vs ...Validator[T]) Validator[*T] {
	return &pointerValidator[*T, T]{
		vs: vs,
	}
}

type pointerValidator[P ~*T, T any] struct {
	vs []Validator[T]
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *pointerValidator[P, T]) WithFormat(key message.Reference, a ...Arg) Validator[P] {
	rr := *r
	//s.format = xxx
	return &rr
}

// Validate validates v.
func (r *pointerValidator[P, T]) Validate(ctx context.Context, p P) error {
	if p == nil {
		return nil
	}
	v := *p
	var errs []error
	for _, rule := range r.vs {
		if err := rule.Validate(ctx, v); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return joinErrors(errs...)
	}
	return nil
}

// pointerError reports an error is caused in Pointer validator.
type pointerError[P ~*T, T any] struct {
	Value P
}
