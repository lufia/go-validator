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
			User    User
			Options []string
		}
	)

	var requestValidator = validator.Struct(func(s validator.StructRuleAdder, r *Request) {
		s.Add(validator.Field(&r.User, "user"), validator.Struct(func(s validator.StructRuleAdder, u *User) {
			s.Add(validator.Field(&u.ID, "id"), validator.Length[string](5, 10))
			s.Add(validator.Field(&u.Name, "name"), validator.Required[string]())
		}))
		s.Add(validator.Field(&r.Options, "options"), validator.Slice[string](
			validator.In("option1", "option2")),
		)
	})

	var r Request
	r.Options = []string{"option3"}
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
	// options: must be a valid value in [option1 option2]
}

func Example_separated() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User    User
			Options []string
		}
	)

	var (
		userValidator = validator.Struct(func(s validator.StructRuleAdder, u *User) {
			s.Add(validator.Field(&u.ID, "id"), validator.Length[string](5, 10))
			s.Add(validator.Field(&u.Name, "name"), validator.Required[string]())
		})
		requestValidator = validator.Struct(func(s validator.StructRuleAdder, r *Request) {
			s.Add(validator.Field(&r.User, "user"), userValidator)
			s.Add(validator.Field(&r.Options, "options"), validator.Slice[string](
				validator.In("option1", "option2")),
			)
		})
	)

	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
}
