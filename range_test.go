package validator

import (
	"testing"
)

func TestMin(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Min(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 4, "")
		testValidate(t, v, 2, "must be no less than 3")
	})
}

func TestMinWithCustomReferenceKey(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Min(3).WithReferenceKey("less than %v", ByName("min"))
		testValidate(t, v, 2, "less than 3")
	})
}

func TestMax(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Max(3)
		testValidate(t, v, 3, "")
		testValidate(t, v, 2, "")
		testValidate(t, v, 4, "must be no greater than 3")
	})
}

func TestMaxWithCustomReferenceKey(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := Max(3).WithReferenceKey("greater than %v", ByName("max"))
		testValidate(t, v, 4, "greater than 3")
	})
}

func TestInRange(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := InRange(1, 3)
		testValidate(t, v, 1, "")
		testValidate(t, v, 3, "")
		testValidate(t, v, 0, "must be in range(1 ... 3)")
		testValidate(t, v, 4, "must be in range(1 ... 3)")
	})
}

func TestInRangeWithCustomReferenceKey(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := InRange(1, 3).WithReferenceKey("out of range(%v, %v)", ByName("min"), ByName("max"))
		testValidate(t, v, 4, "out of range(1, 3)")
	})
}
