package requiring

import (
	"fmt"
	"io"
	"testing"
)

type testRequiredViolationPrinter[T comparable] struct{}

func (testRequiredViolationPrinter[T]) Print(w io.Writer, e RequiredViolationError[T]) {
	fmt.Fprintf(w, "'%v' is empty", e.Value)
}

var _ RequiredViolationPrinter[string] = (*testRequiredViolationPrinter[string])(nil)

func TestRequiredVerified_string(t *testing.T) {
	tests := []string{
		"a",
		"ab",
		"\n",
	}
	for _, s := range tests {
		err := Required[string]().Validate(s)
		if err != nil {
			t.Fatalf("Validate(%q): %v", s, err)
		}
	}
}

func TestRequiredViolation_string(t *testing.T) {
	tests := map[string]struct {
		errstr string
		opts   []any
	}{
		"default": {
			errstr: "required",
		},
		"with_printer": {
			errstr: "'' is empty",
			opts: []any{
				&testRequiredViolationPrinter[string]{},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Required[string](tt.opts...).Validate("")
			if err == nil {
				t.Fatalf(`Validate("") should return a violation error`)
			}
			if s := err.Error(); s != tt.errstr {
				t.Errorf(`Validate("") = %q; want %q`, s, tt.errstr)
			}
		})
	}
}
