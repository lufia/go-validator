package validator

import (
	"testing"
)

func TestPattern(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := PatternString[string]("^ab*$")
		testValidate(t, v, "a", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "", "must match the pattern /^ab*$/")
	})
}

func TestPatternWithFormat(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := PatternString[string]("123").WithFormat("does not match")
		testValidate(t, v, "", "does not match")
	})
}
