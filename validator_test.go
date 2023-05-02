package requiring

import (
	"testing"
)

func Test_NotEmpty_string(t *testing.T) {
	tests := map[string]struct {
		Value string
		Err   string
	}{
		"empty string": {"", "it is required"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := NotEmpty.Validate(tt.Value)
			if err == nil {
				t.Fatalf("NotEmpty.Validate(%q) should return a violation error", tt.Value)
			}
			if s := err.Error(); s != tt.Err {
				t.Errorf("NotEmpty.Validate(%q) = %q; want %q", tt.Value, s, tt.Err)
			}
		})
	}
}
