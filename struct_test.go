package requiring

import (
	"testing"
)

func TestStructVerified(t *testing.T) {
	type (
		Request struct {
			Name string
			Type int
		}
	)
	v := Struct(func(s StructRuleAdder, r *Request) {
		s.Add(&r.Name, "name")
	})
	var r Request
	if err := v.Validate(&r); err != nil {
		t.Errorf("Validate(%#v): %v", r, err)
	}
}
