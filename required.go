package validator

import (
	"context"
	"errors"

	"golang.org/x/text/message"
)

// Required returns the validator to verify the value is not zero value.
//
// A named arg is available in its error format.
//   - value: user input (type T)
func Required[T comparable]() Validator[T] {
	return &requiredValidator[T]{
		format: requiredErrorFormat,
	}
}

// requiredValidator represents the validator to check the value is not zero-value.
type requiredValidator[T comparable] struct {
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *requiredValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *requiredValidator[T]) Validate(ctx context.Context, v T) error {
	var v0 T
	if v == v0 {
		e := &requiredError[T]{
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// requiredError reports an error is caused in Required validator.
type requiredError[T comparable] struct {
	Value T `arg:"value"`
}

var _ Validator[string] = (*requiredValidator[string])(nil)
