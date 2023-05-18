package validator

import (
	"testing"
)

func TestSlice(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Slice(Required[string]())
		testValidate(t, v, []string{"a", "ab"}, "")
		testValidate(t, v, []string(nil), "")
		testValidate(t, v, []string{""}, "cannot be the zero value")
	})
	t.Run("int", func(t *testing.T) {
		v := Slice(Required[int]())
		testValidate(t, v, []int{1, -1}, "")
	})
}
