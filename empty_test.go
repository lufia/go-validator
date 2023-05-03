package requiring

import (
	"fmt"
	"io"
	"testing"
)

type testNotEmptyViolationPrinter[T ~string] struct{}

func (testNotEmptyViolationPrinter[T]) Print(w io.Writer, e NotEmptyViolationError[T]) {
	fmt.Fprintf(w, "'%s' is empty", e.Value)
}

var _ NotEmptyViolationPrinter[string] = (*testNotEmptyViolationPrinter[string])(nil)

func TestNotEmptyVerified_string(t *testing.T) {
	tests := []string{
		"a",
		"ab",
		"\n",
	}
	for _, s := range tests {
		err := NotEmpty[string]().Validate(s)
		if err != nil {
			t.Fatalf("Validate(%q): %v", s, err)
		}
	}
}

func TestNotEmptyViolation_string(t *testing.T) {
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
				&testNotEmptyViolationPrinter[string]{},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := NotEmpty[string](tt.opts...).Validate("")
			if err == nil {
				t.Fatalf(`Validate("") should return a violation error`)
			}
			if s := err.Error(); s != tt.errstr {
				t.Errorf(`Validate("") = %q; want %q`, s, tt.errstr)
			}
		})
	}
}
