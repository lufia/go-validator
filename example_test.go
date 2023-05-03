package requiring_test

import (
	"fmt"

	"github.com/lufia/go-requiring"
)

type Request struct {
	User User
}

type User struct {
	Name string
}

var requestValidator = requiring.Struct(func(s requiring.StructRuleAdder, r *Request) {
	s.Add(&r.User, "user", requiring.Struct(func(s requiring.StructRuleAdder, u *User) {
		s.Add(&u.Name, "name", requiring.NotEmpty[string]())
	}))
})

func Example() {
	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Output: 'name' is required
}
