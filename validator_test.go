package requiring

import (
	"fmt"
	"io"
	"testing"
)

type testInvalidTypePrinter struct{}

func (testInvalidTypePrinter) Print(w io.Writer, e InvalidTypeError) {
	fmt.Fprintf(w, "'%s' %[2]T(%[2]v) vs %[3]v", e.Name, e.Value, e.Type)
}

var _ InvalidTypePrinter = (*testInvalidTypePrinter)(nil)

func testValidate(t *testing.T, v Validator, p any, e string) {
	t.Helper()

	err := v.Validate(p)
	if e == "" {
		if err != nil {
			t.Errorf("Validate(%v) should be passed; but got %v", p, err)
		}
		return
	}

	if err == nil || err.Error() != e {
		t.Errorf("Validate(%v) = %v; want %v", p, err, e)
	}
}
