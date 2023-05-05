package validator

import (
	"bytes"
	"fmt"
)

type mapError map[string]error

func (m mapError) Error() string {
	var w bytes.Buffer
	for key, err := range m {
		fmt.Fprintf(&w, "%s:\n%s", key, err)
	}
	return w.String()
}

// Unwrap returns multiple errors consists in m, like errors.Join.
func (m mapError) Unwrap() []error {
	errs := make([]error, 0, len(m))
	for _, err := range m {
		errs = append(errs, err)
	}
	return errs
}

func unwrapErrors(err error) ([]error, bool) {
	errs, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return nil, false
	}
	return errs.Unwrap(), true
}
