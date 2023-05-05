package validator

import (
	"errors"
	"testing"
)

func TestMapError(t *testing.T) {
	tests := map[string]struct {
		err mapError
		msg string
	}{
		"nil": {
			nil,
			"",
		},
		"one": {
			mapError{"key1": errors.New("err1")},
			"key1:\nerr1",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if e := tt.err.Error(); e != tt.msg {
				t.Errorf("Error() = %q; want %q", e, tt.msg)
			}
		})
	}
}
