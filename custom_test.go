package validator_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lufia/go-validator"
	"golang.org/x/text/message"
)

type CreateUserRequest struct {
	Name                 string
	Password             string
	ConfirmationPassword string
}

type UsernameValidator struct{}

func (*UsernameValidator) Validate(ctx context.Context, v string) error {
	// find non-alnum or non-ascii character
	i := strings.IndexFunc(v, func(c rune) bool {
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

func (r *UsernameValidator) WithFormat(key message.Reference, a ...validator.Arg) validator.Validator[string] {
	rr := *r
	return &rr
}

type CreateUserRequestValidator struct{}

func (*CreateUserRequestValidator) Validate(ctx context.Context, v *CreateUserRequest) error {
	if v.Password != v.ConfirmationPassword {
		return errors.New("passwords does not match")
	}
	return nil
}

func (r *CreateUserRequestValidator) WithFormat(key message.Reference, a ...validator.Arg) validator.Validator[*CreateUserRequest] {
	rr := *r
	return &rr
}

var createUserRequestValidator = validator.Join(
	validator.Struct(func(s validator.StructRule, r *CreateUserRequest) {
		validator.AddField(s, &r.Name, "name",
			validator.Length[string](5, 20),
			validator.Validator[string](&UsernameValidator{}))
		validator.AddField(s, &r.Password, "password",
			validator.MinLength[string](8))
		validator.AddField(s, &r.ConfirmationPassword, "confirmation-password",
			validator.MinLength[string](8))
	}),
	validator.Validator[*CreateUserRequest](&CreateUserRequestValidator{}),
)

func Example_customValidator() {
	ctx := context.Background()
	err := createUserRequestValidator.Validate(ctx, &CreateUserRequest{
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
