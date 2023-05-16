package validator

import (
	"testing"
)

func TestRequired(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Required[string]()
		testValidate(t, v, "a", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "", "cannot be the zero value")
	})
	t.Run("int", func(t *testing.T) {
		v := Required[int]()
		testValidate(t, v, 1, "")
		testValidate(t, v, -1, "")
	})
}

func TestRequiredWithFormat(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Required[string]().WithFormat("is empty")
		testValidate(t, v, "", "is empty")
	})
}
