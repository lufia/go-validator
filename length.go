package validator

import (
	"context"
	"errors"

	"golang.org/x/text/message"
)

// MinLength returns the validator to verify the length of the value is greater or equal than n.
//
// Two named args are available in its error format.
//   - min: specified min value (type int)
//   - value: user input (type T)
func MinLength[T ~string](n int) Validator[T] {
	return &minLengthValidator[T]{
		min:    n,
		format: minLengthErrorFormat,
	}
}

// minLengthValidator represents the validator to check the length of the value is greater or equal than T.
type minLengthValidator[T ~string] struct {
	min    int
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *minLengthValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *minLengthValidator[T]) Validate(ctx context.Context, v T) error {
	a := []rune(v)
	if len(a) < r.min {
		e := &minLengthError[T]{
			Value: v,
			Min:   r.min,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// minLengthError reports an error is caused in MinLength validator.
type minLengthError[T ~string] struct {
	Min   int `arg:"min"`
	Value T   `arg:"value"`
}

var _ Validator[string] = (*minLengthValidator[string])(nil)

// MaxLength returns the validator to verify the length of the value is less or equal than n.
//
// Two named args are available in its error format.
//   - max: specified max value (type int)
//   - value: user input (type T)
func MaxLength[T ~string](n int) Validator[T] {
	return &maxLengthValidator[T]{
		max:    n,
		format: maxLengthErrorFormat,
	}
}

// maxLengthValidator represents the validator to check the length of the value is less or equal than T.
type maxLengthValidator[T ~string] struct {
	max    int
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *maxLengthValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *maxLengthValidator[T]) Validate(ctx context.Context, v T) error {
	a := []rune(v)
	if len(a) > r.max {
		e := &maxLengthError[T]{
			Value: v,
			Max:   r.max,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// maxLengthError reports an error is caused in MaxLength validator.
type maxLengthError[T ~string] struct {
	Max   int `arg:"max"`
	Value T   `arg:"value"`
}

var _ Validator[string] = (*maxLengthValidator[string])(nil)

// Length returns the validator to verify the length of the value is within min and max.
//
// Three named args are available in its error format.
//   - min: specified min value (type int)
//   - max: specified max value (type int)
//   - value: user input (type T)
func Length[T ~string](min, max int) Validator[T] {
	return &lengthValidator[T]{
		min:    min,
		max:    max,
		format: lengthErrorFormat,
	}
}

// lengthValidator represents the validator to check the length of the value is within T.
type lengthValidator[T ~string] struct {
	min    int
	max    int
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *lengthValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *lengthValidator[T]) Validate(ctx context.Context, v T) error {
	a := []rune(v)
	if len(a) < r.min || len(a) > r.max {
		e := &lengthError[T]{
			Min:   r.min,
			Max:   r.max,
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// lengthError reports an error is caused in Length validator.
type lengthError[T ~string] struct {
	Min   int `arg:"min"`
	Max   int `arg:"max"`
	Value T   `arg:"value"`
}

var _ Validator[string] = (*lengthValidator[string])(nil)
