//go:build !go1.20

package validator

import "strings"

type joinError struct {
	errs []error
}

func joinErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	e := make([]error, len(errs))
	for i, err := range errs {
		e[i] = err
	}
	return &joinError{e}
}

func (e *joinError) Error() string {
	a := make([]string, len(e.errs))
	for i, err := range e.errs {
		a[i] = err.Error()
	}
	return strings.Join(a, "\n")
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
