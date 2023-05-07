package validator_test

import (
	"fmt"

	"github.com/lufia/go-validator"
)

func Example() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User User
		}
	)

	var requestValidator = validator.Struct(func(s validator.StructRuleAdder, r *Request) {
		s.Add(validator.Field(&r.User, "user"), validator.Struct(func(s validator.StructRuleAdder, u *User) {
			s.Add(validator.Field(&u.ID, "id"), validator.Length[string](5, 10))
			s.Add(validator.Field(&u.Name, "name"), validator.Required[string]())
		}))
	})

	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
}

func Example_separated() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User User
		}
	)

	var (
		userValidator = validator.Struct(func(s validator.StructRuleAdder, u *User) {
			s.Add(validator.Field(&u.ID, "id"), validator.Length[string](5, 10))
			s.Add(validator.Field(&u.Name, "name"), validator.Required[string]())
		})
		requestValidator = validator.Struct(func(s validator.StructRuleAdder, r *Request) {
			s.Add(validator.Field(&r.User, "user"), userValidator)
		})
	)

	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
}
