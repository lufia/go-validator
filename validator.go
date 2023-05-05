// Package validator provides utilities for validating any types.
package validator

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type ViolationError interface {
	error
}

type Printer[E ViolationError] interface {
	Print(w io.Writer, e E)
}

type Validator interface {
	Validate(v any) error
}

type InvalidTypeError struct {
	Value any          // passed value
	Type  reflect.Type // expected type

	p InvalidTypePrinter
}

func (e InvalidTypeError) Error() string {
	p := e.p
	if p == nil {
		p = &invalidTypePrinter{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type invalidTypePrinter struct{}

func (invalidTypePrinter) Print(w io.Writer, e InvalidTypeError) {
	fmt.Fprintf(w, "the value %v is %T, not %v", e.Value, e.Value, e.Type)
}

type InvalidTypePrinter interface {
	Printer[InvalidTypeError]
}

var (
	_ ViolationError     = (*InvalidTypeError)(nil)
	_ InvalidTypePrinter = (*invalidTypePrinter)(nil)
)
