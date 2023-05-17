# go-validator
Yet another validator written in Go.

[![GoDev][godev-image]][godev-url]
[![Actions Status][actions-image]][actions-url]
[![Coverage Status][coveralls-image]][coveralls-url]

## Features

* strongly-typed validators by type parameters
* i18n support in the standard way
* handling multiple validation errors

## Built-in validators

* **Required**: validates comparable types if the value is not zero-value.
* **Length**: validates strings if the length of the value is within the range.
* **MinLength**: see **Length**.
* **MaxLength**: see **Length**.
* **InRange**: validates ordered types if the value is within the range.
* **Min**: see **InRange**.
* **Max**: see **InRange**.
* **In**: validates comparable types if the value is in valid values.
* **Pattern**: validates strings if it matches the regular expression.

## Supported languages

* English
* Japanese

## Example

```go
import (
	"context"
	"fmt"

	"github.com/lufia/go-validator"
)

type OIDCProvider int

const (
	Google OIDCProvider = iota + 1
	Apple
	GitHub
)

type CreateUserRequest struct {
	Name     string
	Provider OIDCProvider
	Theme    string
}

var createUserRequestValidator = validator.Struct(func(s validator.StructRule, r *CreateUserRequest) {
	validator.AddField(s, &r.Name, "name", validator.Length[string](5, 20))
	validator.AddField(s, &r.Provider, "provider", validator.In(Google, Apple, GitHub))
	validator.AddField(s, &r.Theme, "theme", validator.In("light", "dark"))
})

func main() {
	var r CreateUserRequest
	err := createUserRequestValidator.Validate(context.Background(), &r)
	fmt.Println(err)
}
```

For more details, see [the documentation][godev-url].

[godev-image]: https://pkg.go.dev/badge/github.com/lufia/go-validator
[godev-url]: https://pkg.go.dev/github.com/lufia/go-validator
[actions-image]: https://github.com/lufia/go-validator/workflows/Test/badge.svg?branch=main
[actions-url]: https://github.com/lufia/go-validator/actions?workflow=Test
[coveralls-image]: https://coveralls.io/repos/github/lufia/go-validator/badge.svg
[coveralls-url]: https://coveralls.io/github/lufia/go-validator
