package validator

import (
	"testing"
)

func TestMinLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MinLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "1234", "")
		testValidate(t, v, "ab", "the length must be no less than 3")
	})
}

func TestMinLengthWithCustomReferenceKey(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MinLength[string](3).WithReferenceKey("less than %v", ByName("min"))
		testValidate(t, v, "ab", "less than 3")
	})
}

func TestMaxLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MaxLength[string](3)
		testValidate(t, v, "abc", "")
		testValidate(t, v, "ab", "")
		testValidate(t, v, "1234", "the length must be no greater than 3")
	})
}

func TestMaxLengthWithCustomReferenceKey(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := MaxLength[string](3).WithReferenceKey("greater than %v", ByName("max"))
		testValidate(t, v, "1234", "greater than 3")
	})
}

func TestLength(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Length[string](1, 3)
		testValidate(t, v, "a", "")
		testValidate(t, v, "abc", "")
		testValidate(t, v, "", "the length must be in range(1 ... 3)")
		testValidate(t, v, "1234", "the length must be in range(1 ... 3)")
	})
}

func TestLengthWithCustomReferenceKey(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Length[string](1, 3).WithReferenceKey("out of range(%v, %v)", ByName("min"), ByName("max"))
		testValidate(t, v, "1234", "out of range(1, 3)")
	})
}
