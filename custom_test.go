package validator_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/lufia/go-validator"
)

type CreateUserRequest struct {
	Name                 string
	Password             string
	ConfirmationPassword string
}

var (
	usernameValidator = validator.New(func(ctx context.Context, v string) bool {
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
		return i < 0
	}).WithFormat("does not allow not-alphabets or not-digits")

	passwordValidator = validator.New(func(ctx context.Context, r *CreateUserRequest) bool {
		return r.Password == r.ConfirmationPassword
	}).WithFormat("passwords does not match")
)

var createUserRequestValidator = validator.Join(
	validator.Struct(func(s validator.StructRule, r *CreateUserRequest) {
		validator.AddField(s, &r.Name, "name",
			validator.Length[string](5, 20),
			usernameValidator)
		validator.AddField(s, &r.Password, "password",
			validator.MinLength[string](8))
		validator.AddField(s, &r.ConfirmationPassword, "confirmation-password",
			validator.MinLength[string](8))
	}),
	passwordValidator,
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
