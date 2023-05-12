package validator

import (
	"context"
	"errors"

	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	minKey     = "must be no less than %v"
	maxKey     = "must be no greater than %v"
	inRangeKey = "must be in range(%v ... %v)"
)

func init() {
	DefaultCatalog.SetString(language.English, minKey, minKey)
	DefaultCatalog.SetString(language.Japanese, minKey, "xxx")

	DefaultCatalog.SetString(language.English, maxKey, maxKey)
	DefaultCatalog.SetString(language.Japanese, maxKey, "xxx")

	DefaultCatalog.SetString(language.English, inRangeKey, inRangeKey)
	DefaultCatalog.SetString(language.Japanese, inRangeKey, "xxx")
}

type ordered interface {
	constraints.Ordered
}

// Min returns the validator to verify the value is greater or equal than n.
//
// This validator has an arg in its reference key.
//   - min: specified min value (type T)
//   - value: user input (type T)
func Min[T ordered](n T) Validator[T] {
	return &minValidator[T]{
		min:  n,
		key:  minKey,
		args: []Arg{ByName("min")},
	}
}

// minValidator represents the validator to check the value is greater or equal than T.
type minValidator[T ordered] struct {
	min  T
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *minValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
	return &rr
}

// Validate validates v.
func (r *minValidator[T]) Validate(ctx context.Context, v T) error {
	if v < r.min {
		e := &minError[T]{
			Min:   r.min,
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
// This validator has an arg in its reference key.
//   - max: specified max value (type T)
//   - value: user input (type T)
func Max[T ordered](n T) Validator[T] {
	return &maxValidator[T]{
		max:  n,
		key:  maxKey,
		args: []Arg{ByName("max")},
	}
}

// maxValidator represents the validator to check the value is less or equal than T.
type maxValidator[T ordered] struct {
	max  T
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *maxValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
	return &rr
}

// Validate validates v.
func (r *maxValidator[T]) Validate(ctx context.Context, v T) error {
	if v > r.max {
		e := &maxError[T]{
			Value: v,
			Max:   r.max,
		}
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
// This validator has two args in its reference key.
//   - min: specified min value (type T)
//   - max: specified max value (type T)
//   - value: user input (type T)
func InRange[T ordered](min, max T) Validator[T] {
	return &inRangeValidator[T]{
		min:  min,
		max:  max,
		key:  inRangeKey,
		args: []Arg{ByName("min"), ByName("max")},
	}
}

// inRangeValidator represents the validator to check the value is within T.
type inRangeValidator[T ordered] struct {
	min  T
	max  T
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *inRangeValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
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
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
