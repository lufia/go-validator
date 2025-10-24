package validator

import (
	"context"
	"reflect"
	"testing"
)

func testValidate[V Validator[T], T any](t *testing.T, v V, p T, e string) {
	t.Helper()

	err := v.Validate(context.Background(), p)
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
		vs  []Validator[int]
	}{
		"none":   {0, nil},
		"simple": {1, []Validator[int]{Min(1)}},
		"multi":  {2, []Validator[int]{Min(1), Max(10)}},
		"nested": {3, []Validator[int]{Required[int](), Join(Min(1), Max(10))}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v := Join(tt.vs...).(*joinValidator[int])
			if n := len(v.vs); n != tt.num {
				t.Errorf("got %d; want %d", n, tt.num)
			}
			if err := v.Validate(context.Background(), 2); err != nil {
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
		err := v.Validate(context.Background(), s)
		if err != nil {
			t.Errorf("Validate(%q) = %v", s, err)
		}
	})
	t.Run("error", func(t *testing.T) {
		s := ""
		err := v.Validate(context.Background(), s)
		if err == nil {
			t.Errorf("Validate(%q) should return an error", s)
		}
	})
}

func TestOrderedMap(t *testing.T) {
	tests := map[string][]string{
		"names":  {"sys", "dev", "bin"},
		"digits": {"1", "2", "3"},
	}
	for name, a := range tests {
		t.Run(name, func(t *testing.T) {
			var m OrderedMap[string, int]
			for i, s := range a {
				m.set(s, i)
			}
			keys := m.Keys()
			if !reflect.DeepEqual(keys, a) {
				t.Errorf("Keys() = %v; want %v", keys, a)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := New(func(ctx context.Context, s string) bool { return s != "" })
		testValidate(t, v, "a", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "", "must be a valid value")
	})
	t.Run("int", func(t *testing.T) {
		v := New(func(ctx context.Context, n int) bool { return n >= 0 })
		testValidate(t, v, 1, "")
		testValidate(t, v, -1, "must be a valid value")
	})
}

func TestNewWithFormat(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := New(func(ctx context.Context, s string) bool { return s != "" }).WithFormat("is empty")
		testValidate(t, v, "", "is empty")
	})
}
