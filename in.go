package validator

import (
	"context"
	"errors"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const inKey = "must be a valid value in %v"

func init() {
	DefaultCatalog.SetString(language.English, inKey, "must be a valid value in %v")
	DefaultCatalog.SetString(language.Japanese, inKey, "xxxx")
}

// In returns the validator to verify the value is in a.
//
// This validator has two args in its reference key.
//   - validValues: specified valid values (type []T)
//   - value: user input (type T)
func In[T comparable](a ...T) Validator[T] {
	return &inValidator[T]{
		a:    a,
		key:  inKey,
		args: []Arg{ByName("validValues")},
	}
}

// inValidator represents the validator to check the value is in T.
type inValidator[T comparable] struct {
	a    []T
	key  message.Reference
	args []Arg
}

// WithReferenceKey returns shallow copy of r with its reference key changed to key.
func (r *inValidator[T]) WithReferenceKey(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.key = key
	rr.args = a
	return &rr
}

// Validate validates v.
func (r *inValidator[T]) Validate(ctx context.Context, v T) error {
	if !slices.Contains(r.a, v) {
		e := &inError[T]{
			Value:       v,
			ValidValues: r.a,
		}
		return errors.New(ctxPrint(ctx, e, r.key, r.args))
	}
	return nil
}

// inError reports an error is caused in In validator.
type inError[T comparable] struct {
	Value       T   `arg:"value"`
	ValidValues []T `arg:"validValues"`
}

var _ Validator[string] = (*inValidator[string])(nil)
