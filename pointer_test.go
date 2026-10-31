package validator

import (
	"testing"
)

func TestPointer(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Pointer(Required[string]())
		testValidate(t, v, new("test"), "")
		testValidate(t, v, new(""), "cannot be the zero value")
		testValidate(t, v, nil, "")
	})
}
