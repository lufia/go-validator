# go-validator
Yet another validator written in Go.

[![GoDev][godev-image]][godev-url]
[![Actions Status][actions-image]][actions-url]
[![Coverage Status][coveralls-image]][coveralls-url]

## Description

TODO

## Built-in validators

* **Required**: validates comparable types if the value is not zero-value.
* **Length**: validates strings if the length of the value is within the range.
* **MinLength**: see **Length**.
* **MaxLength**: see **Length**.
* **InRange**: validates ordered types if the value is within the range.
* **Min**: see **InRange**.
* **Max**: see **InRange**.
* **In**: validates comparable types if the value is in valid values.

[godev-image]: https://pkg.go.dev/badge/github.com/lufia/go-validator
[godev-url]: https://pkg.go.dev/github.com/lufia/go-validator
[actions-image]: https://github.com/lufia/go-validator/workflows/Test/badge.svg?branch=main
[actions-url]: https://github.com/lufia/go-validator/actions?workflow=Test
[coveralls-image]: https://coveralls.io/repos/github/lufia/go-validator/badge.svg
[coveralls-url]: https://coveralls.io/github/lufia/go-validator
