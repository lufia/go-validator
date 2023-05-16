package validator

import (
	"testing"
)

func TestIn(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		v := In("a", "b")
		testValidate(t, v, "a", "")
		testValidate(t, v, "b", "")
		testValidate(t, v, "x", "must be a valid value in [a b]")
	})
	t.Run("int", func(t *testing.T) {
		v := In(1, 2)
		testValidate(t, v, 1, "")
		testValidate(t, v, 2, "")
	})
}

func TestInWithFormat(t *testing.T) {
	t.Run("printer", func(t *testing.T) {
		v := In("a", "b").WithFormat("must in %v", ByName("validValues"))
		testValidate(t, v, "x", "must in [a b]")
	})
}
