package validator

import (
	"context"
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	minLengthKey = "the length must be no less than %v"
	maxLengthKey = "the length must be no greater than %v"
	lengthKey    = "the length must be in range(%v ... %v)"
)

func init() {
	DefaultCatalog.SetString(language.English, minLengthKey, minLengthKey)
	DefaultCatalog.SetString(language.Japanese, minLengthKey, "xxx")

	DefaultCatalog.SetString(language.English, maxLengthKey, maxLengthKey)
	DefaultCatalog.SetString(language.Japanese, maxLengthKey, "xxx")

	DefaultCatalog.SetString(language.English, lengthKey, lengthKey)
	DefaultCatalog.SetString(language.Japanese, lengthKey, "xxx")
}

// MinLength returns the validator to verify the length of the value is greater or equal than n.
//
// This validator has an arg in its reference key.
//   - min: specified min value (type int)
//   - value: user input (type T)
func MinLength[T ~string](n int) Validator[T] {
	return &minLengthValidator[T]{
		min:  n,
		key:  minLengthKey,
		args: []Arg{ByName("min")},
	}
}

// minLengthValidator represents the validator to check the length of the value is greater or equal than T.
type minLengthValidator[T ~string] struct {
	min  int
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *minLengthValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
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
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
// This validator has an arg in its reference key.
//   - max: specified max value (type int)
//   - value: user input (type T)
func MaxLength[T ~string](n int) Validator[T] {
	return &maxLengthValidator[T]{
		max:  n,
		key:  maxLengthKey,
		args: []Arg{ByName("max")},
	}
}

// maxLengthValidator represents the validator to check the length of the value is less or equal than T.
type maxLengthValidator[T ~string] struct {
	max  int
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *maxLengthValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
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
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
// This validator has two args in its reference key.
//   - min: specified min value (type int)
//   - max: specified max value (type int)
//   - value: user input (type T)
func Length[T ~string](min, max int) Validator[T] {
	return &lengthValidator[T]{
		min:  min,
		max:  max,
		key:  lengthKey,
		args: []Arg{ByName("min"), ByName("max")},
	}
}

// lengthValidator represents the validator to check the length of the value is within T.
type lengthValidator[T ~string] struct {
	min  int
	max  int
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *lengthValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
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
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
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
