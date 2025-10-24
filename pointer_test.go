package validator_test

import (
	"context"
	"testing"

	"golang.org/x/text/message"

	"github.com/lufia/go-validator"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    validator.Validator[*int]
		v    *int
		err  string
	}{
		{
			name: "nil",
			r:    validator.Pointer[int](),
			v:    nil,
			err:  "must not be nil",
		},
		{
			name: "not nil",
			r:    validator.Pointer[int](),
			v:    new(int),
			err:  "",
		},
		{
			name: "WithFormat",
			r:    validator.Pointer[int]().WithFormat(message.Reference("pointer error")),
			v:    nil,
			err:  "pointer error",
		},
		{
			name: "nested validation ok",
			r:    validator.Pointer[int](validator.Min(0)),
			v:    ptr(10),
			err:  "",
		},
		{
			name: "nested validation error",
			r:    validator.Pointer[int](validator.Min(100)),
			v:    ptr(10),
			err:  "must be no less than 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Validate(context.Background(), tt.v)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("want error %q, got nil", tt.err)
				}
				if err.Error() != tt.err {
					t.Errorf("want error %q, got %q", tt.err, err.Error())
				}
			} else if err != nil {
				t.Fatalf("want no error, got %q", err.Error())
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestPointer_zeroValue(t *testing.T) {
	t.Parallel()

	var v *int
	r := validator.Pointer[int]()

	err := r.Validate(context.Background(), v)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if err.Error() != "must not be nil" {
		t.Errorf("want %q, got %q", "must not be nil", err.Error())
	}
}

func TestPointer_nonZeroValue(t *testing.T) {
	t.Parallel()

	var v int = 1
	r := validator.Pointer[int]()

	err := r.Validate(context.Background(), &v)
	if err != nil {
		t.Fatalf("want no error, got %q", err.Error())
	}
}

func TestPointer_nestedMinValidator(t *testing.T) {
	t.Parallel()

	val := 5
	r := validator.Pointer[int](validator.Min(10))

	err := r.Validate(context.Background(), &val)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if err.Error() != "must be no less than 10" {
		t.Errorf("want %q, got %q", "must be no less than 10", err.Error())
	}
}

func TestPointer_nestedMinValidator_ok(t *testing.T) {
	t.Parallel()

	val := 15
	r := validator.Pointer[int](validator.Min(10))

	err := r.Validate(context.Background(), &val)
	if err != nil {
		t.Fatalf("want no error, got %q", err.Error())
	}
}

func TestPointer_nestedJoinValidator(t *testing.T) {
	t.Parallel()

	val := -5
	r := validator.Pointer[int](
		validator.Join(
			validator.Min(0),
			validator.Max(10),
		),
	)

	err := r.Validate(context.Background(), &val)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	expectedErrors := []string{
		"must be no less than 0",
	}

	unwrapErr, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("error does not implement Unwrap() []error: %T", err)
	}

	errList := unwrapErr.Unwrap()
	if len(errList) != len(expectedErrors) {
		t.Errorf("expected %d errors, got %d", len(expectedErrors), len(errList))
	}

	for i, e := range errList {
		if e.Error() != expectedErrors[i] {
			t.Errorf("expected error %q, got %q", expectedErrors[i], e.Error())
		}
	}
}

func TestPointer_multipleNestedValidators(t *testing.T) {
	t.Parallel()

	val := 5
	r := validator.Pointer[int](
		validator.Min(0),
		validator.Max(100),
	)

	err := r.Validate(context.Background(), &val)
	if err != nil {
		t.Fatalf("want no error, got %q", err.Error())
	}
}

func TestPointer_multipleNestedValidators_error(t *testing.T) {
	t.Parallel()

	val := 10
	r := validator.Pointer[int](
		validator.Min(0),
		validator.Max(5),
	)

	err := r.Validate(context.Background(), &val)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if err.Error() != "must be no greater than 5" {
		t.Errorf("want %q, got %q", "must be no greater than 5", err.Error())
	}
}

func TestPointer_emptyRules(t *testing.T) {
	t.Parallel()

	val := 123
	r := validator.Pointer[int]()

	err := r.Validate(context.Background(), &val)
	if err != nil {
		t.Fatalf("want no error, got %q", err.Error())
	}
}

func TestPointer_customMessageReference(t *testing.T) {
	t.Parallel()

	r := validator.Pointer[int]().WithFormat(message.Reference("Custom Pointer Error"))

	err := r.Validate(context.Background(), nil)

	if err == nil {
		t.Fatal("want error, got nil")
	}

	if err.Error() != "Custom Pointer Error" {
		t.Errorf("want %q, got %q", "Custom Pointer Error", err.Error())
	}
}

func TestPointer_CustomMessageWithNestedError(t *testing.T) {
	t.Parallel()

	val := -1
	r := validator.Pointer[int](validator.Min(0)).WithFormat(message.Reference("Pointer wrapper error"))

	err := r.Validate(context.Background(), &val)

	if err == nil {
		t.Fatal("want error, got nil")
	}

	joinErr, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected error to be a join error, got %T", err)
	}

	errList := joinErr.Unwrap()
	if len(errList) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errList))
	}

	if errList[0].Error() != "must be no less than 0" {
		t.Errorf("expected nested error %q, got %q", "must be no less than 0", errList[0].Error())
	}
}

func TestPointer_NilNestedErrorCheck(t *testing.T) {
	t.Parallel()

	r := validator.Pointer[int](validator.Required[int]()).WithFormat(message.Reference("Pointer error"))

	err := r.Validate(context.Background(), nil)

	if err == nil {
		t.Fatal("want error for nil pointer, got nil")
	}

	if err.Error() != "Pointer error" {
		t.Errorf("expected error %q, got %q", "Pointer error", err.Error())
	}
}

func TestPointer_NestedErrorWithArgs(t *testing.T) {
	t.Parallel()

	val := 5
	r := validator.Pointer[int](validator.Min(10)).WithFormat(message.Reference("Pointer error: %v"), validator.ByName("value"))

	err := r.Validate(context.Background(), &val)

	if err == nil {
		t.Fatal("want error for nested validation, got nil")
	}

	joinErr, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected error to be a join error, got %T", err)
	}

	errList := joinErr.Unwrap()
	if len(errList) != 1 {
		t.Fatalf("expected 1 nested error, got %d", len(errList))
	}

	if errList[0].Error() != "must be no less than 10" {
		t.Errorf("expected nested error %q, got %q", "must be no less than 10", errList[0].Error())
	}
}

func TestPointer_NilWithNestedRules(t *testing.T) {
	t.Parallel()

	r := validator.Pointer[int](validator.Min(10))

	err := r.Validate(context.Background(), nil)

	if err == nil {
		t.Fatal("want error for nil pointer with nested rules, got nil")
	}

	if err.Error() != "must not be nil" {
		t.Errorf("expected error %q, got %q", "must not be nil", err.Error())
	}
}
