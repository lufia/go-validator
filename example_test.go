package validator_test

import (
	"context"
	"fmt"

	"github.com/lufia/go-validator"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Example() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User    *User
			Options []string
		}
	)

	var requestValidator = validator.Struct(func(s validator.StructRule, r *Request) {
		validator.AddField(s, &r.User, "user", validator.Struct(func(s validator.StructRule, u *User) {
			validator.AddField(s, &u.ID, "id", validator.Length[string](5, 10))
			validator.AddField(s, &u.Name, "name", validator.Required[string]())
		}))
		validator.AddField(s, &r.Options, "options",
			validator.Slice[string](validator.In("option1", "option2")))
	})

	var r Request
	r.Options = []string{"option3"}
	err := requestValidator.Validate(context.Background(), &r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
	// options: must be a valid value in [option1 option2]
}

func Example_localized() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User    *User
			Options []string
		}
	)

	var requestValidator = validator.Struct(func(s validator.StructRule, r *Request) {
		validator.AddField(s, &r.User, "user", validator.Struct(func(s validator.StructRule, u *User) {
			validator.AddField(s, &u.ID, "id", validator.Length[string](5, 10))
			validator.AddField(s, &u.Name, "name", validator.Required[string]())
		}))
		validator.AddField(s, &r.Options, "options",
			validator.Slice[string](validator.In("option1", "option2")))
	})

	p := message.NewPrinter(language.Japanese, message.Catalog(validator.DefaultCatalog))
	ctx := validator.WithPrinter(context.Background(), p)

	var r Request
	r.Options = []string{"option3"}
	err := requestValidator.Validate(ctx, &r)
	fmt.Println(err)
	// Unordered output:
	// user: name: 必須です
	// user: id: 長さは5以上10以内の制限があります
	// options: [option1 option2]のいずれかでなければなりません
}

func Example_separated() {
	type (
		User struct {
			ID   string
			Name string
		}
		Request struct {
			User    *User
			Options []string
		}
	)

	var (
		userValidator = validator.Struct(func(s validator.StructRule, u *User) {
			validator.AddField(s, &u.ID, "id", validator.Length[string](5, 10))
			validator.AddField(s, &u.Name, "name", validator.Required[string]())
		})
		requestValidator = validator.Struct(func(s validator.StructRule, r *Request) {
			validator.AddField(s, &r.User, "user", userValidator)
			validator.AddField(s, &r.Options, "options", validator.Slice[string](
				validator.In("option1", "option2"),
			))
			validator.AddField(s, &r.Options, "options", validator.Slice(
				validator.In("option1", "option2"),
			))
		})
	)

	var r Request
	err := requestValidator.Validate(context.Background(), &r)
	fmt.Println(err)
	// Unordered output:
	// user: name: cannot be the zero value
	// user: id: the length must be in range(5 ... 10)
}
