package validator_test

import (
	"errors"
	"fmt"

	"github.com/lufia/go-validator"
)

type CreateUserRequest struct {
	Name                 string
	Password             string
	ConfirmationPassword string
}

type CreateUserRequestValidator struct{}

func (*CreateUserRequestValidator) Validate(v any) error {
	u := v.(*CreateUserRequest)
	if u.Password != u.ConfirmationPassword {
		return errors.New("passwords does not match")
	}
	return nil
}

var createUserRequestValidator = validator.Join(
	validator.Struct(func(s validator.StructRuleAdder, r *CreateUserRequest) {
		s.Add(
			validator.Field(&r.Name, "name"),
			validator.Length[string](5, 20),
		)
		s.Add(
			validator.Field(&r.Password, "password"),
			validator.MinLength[string](8),
		)
		s.Add(
			validator.Field(&r.ConfirmationPassword, "confirmation-password"),
			validator.MinLength[string](8),
		)
	}),
	&CreateUserRequestValidator{},
)

func Example_customValidator() {
	err := createUserRequestValidator.Validate(&CreateUserRequest{
		Name:                 "test",
		Password:             "1234",
		ConfirmationPassword: "abcd",
	})
	fmt.Println(err)
	// Unordered output:
	// name: the length must be in range(5 ... 20)
	// password: the length must be no less than 8
	// confirmation-password: the length must be no less than 8
	// passwords does not match
}
