package validator

import (
	"cmp"
	"context"
	"errors"

	"golang.org/x/text/message"
)

type ordered interface {
	cmp.Ordered
}

// Min returns the validator to verify the value is greater or equal than n.
//
// Two named args are available in its error format.
//   - min: specified min value (type T)
//   - value: user input (type T)
func Min[T ordered](n T) Validator[T] {
	return &minValidator[T]{
		min:    n,
		format: minErrorFormat,
	}
}

// minValidator represents the validator to check the value is greater or equal than T.
type minValidator[T ordered] struct {
	min    T
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *minValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *minValidator[T]) Validate(ctx context.Context, v T) error {
	if v < r.min {
		e := &minError[T]{
			Min:   r.min,
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// minError reports an error is caused in Min validator.
type minError[T ordered] struct {
	Min   T `arg:"min"`
	Value T `arg:"value"`
}

var _ Validator[int] = (*minValidator[int])(nil)

// Max returns the validator to verify the value is less or equal than n.
//
// Two named args are available in its error format.
//   - max: specified max value (type T)
//   - value: user input (type T)
func Max[T ordered](n T) Validator[T] {
	return &maxValidator[T]{
		max:    n,
		format: maxErrorFormat,
	}
}

// maxValidator represents the validator to check the value is less or equal than T.
type maxValidator[T ordered] struct {
	max    T
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *maxValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *maxValidator[T]) Validate(ctx context.Context, v T) error {
	if v > r.max {
		e := &maxError[T]{
			Value: v,
			Max:   r.max,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// maxError reports an error is caused in Max validator.
type maxError[T ordered] struct {
	Max   T `arg:"max"`
	Value T `arg:"value"`
}

var _ Validator[int] = (*maxValidator[int])(nil)

// InRange returns the validator to verify the value is within min and max.
//
// Three named args are available in its error format.
//   - min: specified min value (type T)
//   - max: specified max value (type T)
//   - value: user input (type T)
func InRange[T ordered](min, max T) Validator[T] {
	return &inRangeValidator[T]{
		min:    min,
		max:    max,
		format: inRangeErrorFormat,
	}
}

// inRangeValidator represents the validator to check the value is within T.
type inRangeValidator[T ordered] struct {
	min    T
	max    T
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *inRangeValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *inRangeValidator[T]) Validate(ctx context.Context, v T) error {
	if v < r.min || v > r.max {
		e := &inRangeError[T]{
			Min:   r.min,
			Max:   r.max,
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// inRangeError reports an error is caused in InRange validator.
type inRangeError[T ordered] struct {
	Min   T `arg:"min"`
	Max   T `arg:"max"`
	Value T `arg:"value"`
}

var _ Validator[int] = (*inRangeValidator[int])(nil)
