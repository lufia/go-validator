package validator

import (
	"context"
	"errors"

	"golang.org/x/exp/slices"
	"golang.org/x/text/message"
)

// In returns the validator to verify the value is in a.
//
// Two named args are available in its error format.
//   - validValues: specified valid values (type []T)
//   - value: user input (type T)
func In[T comparable](a ...T) Validator[T] {
	return &inValidator[T]{
		a:      a,
		format: inErrorFormat,
	}
}

// inValidator represents the validator to check the value is in T.
type inValidator[T comparable] struct {
	a      []T
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *inValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *inValidator[T]) Validate(ctx context.Context, v T) error {
	if !slices.Contains(r.a, v) {
		e := &inError[T]{
			Value:       v,
			ValidValues: r.a,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// inError reports an error is caused in In validator.
type inError[T comparable] struct {
	Value       T   `arg:"value"`
	ValidValues []T `arg:"validValues"`
}

var _ Validator[string] = (*inValidator[string])(nil)
