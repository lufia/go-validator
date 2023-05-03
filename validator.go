// Package requiring provides utilities for validating any types.
package requiring

import (
	"io"
)

type ViolationError interface {
	error
}

type Printer[E ViolationError] interface {
	Print(w io.Writer, e E)
}

type Validator interface {
	SetName(name string)
	Validate(v any) error
}
