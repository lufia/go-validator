package validator

import (
	"testing"

	"github.com/lufia/go-pointer"
)

func TestPointer(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := Pointer(Required[string]())
		testValidate(t, v, pointer.New("test"), "")
		testValidate(t, v, pointer.New(""), "cannot be the zero value")
		testValidate(t, v, nil, "")
	})
}
