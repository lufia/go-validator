package validator_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lufia/go-validator"
)

type CreateUserRequest struct {
	Name                 string
	Password             string
	ConfirmationPassword string
}

type UsernameValidator struct{}

func (*UsernameValidator) Validate(v any) error {
	s := v.(string)
	// find non-alnum or non-ascii character
	i := strings.IndexFunc(s, func(c rune) bool {
		switch {
		default:
			return true
		case c >= 'a' && c <= 'z':
			return false
		case c >= 'A' && c <= 'Z':
			return false
		case c >= '0' && c <= '9':
			return false
		}
	})
	if i >= 0 {
		return errors.New("does not allow not-alphabets or not-digits")
	}
	return nil
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
			&UsernameValidator{},
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
		Name:                 ".adm",
		Password:             "1234",
		ConfirmationPassword: "abcd",
	})
	fmt.Println(err)
	// Unordered output:
	// name: the length must be in range(5 ... 20)
	// name: does not allow not-alphabets or not-digits
	// password: the length must be no less than 8
	// confirmation-password: the length must be no less than 8
	// passwords does not match
}
