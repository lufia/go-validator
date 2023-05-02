package requiring_test

import (
	"fmt"

	"github.com/lufia/go-requiring"
)

type Request struct {
	User string
}

var requestValidator = requiring.Struct(func(s *requiring.RuleSet, r *Request) {
	s.Add(&r.User, "user", requiring.NotEmpty)
})

func Example() {
	var r Request
	err := requestValidator.Validate(&r)
	fmt.Println(err)
	// Output: 'user' is required
}
