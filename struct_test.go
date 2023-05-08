package validator

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestStruct(t *testing.T) {
	type (
		User struct {
			Name string
		}
		Request struct {
			User *User
			Name string
			Type int
		}
	)
	v := Struct(func(s StructRuleAdder, r *Request) {
		s.Add(Field(&r.Name, "name"))
	})
	var r Request
	if err := v.Validate(&r); err != nil {
		t.Errorf("Validate(%#v): %v", r, err)
	}
}

func TestStruct_reusingValidator(t *testing.T) {
	type (
		Request struct {
			Name string
			Key  string
		}
	)
	var r Request
	stringValidator := Required[string]()
	v := Struct(func(s StructRuleAdder, r *Request) {
		s.Add(Field(&r.Name, "name"), stringValidator)
		s.Add(Field(&r.Key, "key"), stringValidator)
	})
	err := v.Validate(&r)
	testErrors[Request](t, err, []string{
		"name: cannot be the zero value",
		"key: cannot be the zero value",
	})
}

func testErrors[T any](t *testing.T, err error, want []string) {
	t.Helper()
	if err == nil {
		t.Errorf("got <nil>; want %#v", want)
		return
	}
	errs := err.(*StructRuleError[T]).Errors
	e := make([]string, 0, len(errs))
	for _, err := range errs {
		e = append(e, err.Error())
	}
	a := slices.Clone(want)
	slices.Sort(a)
	slices.Sort(e)
	if !slices.Equal(e, a) {
		t.Errorf("got %#v; want %#v", e, want)
	}
}
