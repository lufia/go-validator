package validator_test

import (
	"fmt"

	"github.com/lufia/go-validator"
)

type Request struct {
	User User
}

type User struct {
	ID   string
	Name string
}

var requestValidator = validator.Struct(func(s validator.StructRuleAdder, r *Request) {
	s.Add(validator.Field(&r.User, "user"), validator.Struct(func(s validator.StructRuleAdder, u *User) {
		s.Add(validator.Field(&u.ID, "id"), validator.Length[string](5, 10))
		s.Add(validator.Field(&u.Name, "name"), validator.Required[string]())
	}))
})

func Example() {
	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// name: cannot be the zero value
	// id: must be in range(5 ... 10)
}
