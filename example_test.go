package requiring_test

import (
	"fmt"

	"github.com/lufia/go-requiring"
)

type Request struct {
	User User
}

type User struct {
	ID   string
	Name string
}

var requestValidator = requiring.Struct(func(s requiring.StructRuleAdder, r *Request) {
	s.Add(&r.User, "user", requiring.Struct(func(s requiring.StructRuleAdder, u *User) {
		s.Add(&u.ID, "id", requiring.Length[string](5, 10))
		s.Add(&u.Name, "name", requiring.Required[string]())
	}))
})

func Example() {
	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Unordered output:
	// the field 'name' cannot be the zero value
	// the length of the field 'id' must be in range(5 ... 10)
}
