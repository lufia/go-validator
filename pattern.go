package validator

import (
	"context"
	"errors"
	"regexp"

	"golang.org/x/text/message"
)

// Pattern returns the validator to verify the value matches re.
//
// Two named args are available in its error format.
//   - pattern: specific regular expression (type *regexp.Regexp)
//   - value: user input (type T)
func Pattern[T ~string](re *regexp.Regexp) Validator[T] {
	return &patternValidator[T]{
		re:     re,
		format: patternErrorFormat,
	}
}

// PatternString returns the validator to verify the value matches pattern.
//
// When pattern is invalid, PatternString panics.
func PatternString[T ~string](pattern string) Validator[T] {
	re := regexp.MustCompile(pattern)
	return Pattern[T](re)
}

// patternValidator represents the validator to check the value matches the pattern.
type patternValidator[T ~string] struct {
	re     *regexp.Regexp
	format *errorFormat
}

// WithFormat returns shallow copy of r with its error format changed to key.
func (r *patternValidator[T]) WithFormat(key message.Reference, a ...Arg) Validator[T] {
	rr := *r
	rr.format = &errorFormat{Key: key, Args: a}
	return &rr
}

// Validate validates v.
func (r *patternValidator[T]) Validate(ctx context.Context, v T) error {
	if !r.re.MatchString(string(v)) {
		e := &patternError[T]{
			Pattern: r.re,
			Value:   v,
		}
		return errors.New(ctxPrint(ctx, e, r.format.Key, r.format.Args))
	}
	return nil
}

// patternError reports an error is caused in Pattern validator.
type patternError[T ~string] struct {
	Pattern *regexp.Regexp `arg:"pattern"`
	Value   T              `arg:"value"`
}

var _ Validator[string] = (*patternValidator[string])(nil)
