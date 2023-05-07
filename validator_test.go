package validator

import (
	"testing"
)

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

func TestJoin(t *testing.T) {
	tests := map[string]struct {
		num int
		vs  []Validator
	}{
		"none":   {0, nil},
		"simple": {1, []Validator{Min(1)}},
		"multi":  {2, []Validator{Min(1), Max(10)}},
		"nested": {3, []Validator{Required[int](), Join(Min(1), Max(10))}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v := Join(tt.vs...).(*joinValidator)
			if n := len(v.vs); n != tt.num {
				t.Errorf("got %d; want %d", n, tt.num)
			}
			if err := v.Validate(2); err != nil {
				t.Errorf("Validate(2) = %v; want <nil>", err)
			}
		})
	}
}

func TestJoinValidator_Validate(t *testing.T) {
	v := Join(
		Required[string](),
		MinLength[string](1),
	)
	t.Run("passed", func(t *testing.T) {
		s := "a"
		err := v.Validate(s)
		if err != nil {
			t.Errorf("Validate(%q) = %v", s, err)
		}
	})
	t.Run("error", func(t *testing.T) {
		s := ""
		err := v.Validate(s)
		if err == nil {
			t.Errorf("Validate(%q) should return an error", s)
		}
	})
}
