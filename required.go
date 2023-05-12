package validator

import (
	"context"
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const requiredKey = "cannot be the zero value"

func init() {
	DefaultCatalog.SetString(language.English, requiredKey, "cannot be the zero value")
	DefaultCatalog.SetString(language.Japanese, requiredKey, "必須です")
}

// Required returns the validator to verify the value is not zero value.
//
// This validator has an args in its reference key.
//   - value: user input (type T)
func Required[T comparable]() Validator[T] {
	return &requiredValidator[T]{
		key: requiredKey,
	}
}

// requiredValidator represents the validator to check the value is not zero-value.
type requiredValidator[T comparable] struct {
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *requiredValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
	return &rr
}

// Validate validates v.
func (r *requiredValidator[T]) Validate(ctx context.Context, v T) error {
	var v0 T
	if v == v0 {
		e := &requiredError[T]{
			Value: v,
		}
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
	}
	return nil
}

// requiredError reports an error is caused in Required validator.
type requiredError[T comparable] struct {
	Value T `arg:"value"`
}

var _ Validator[string] = (*requiredValidator[string])(nil)
