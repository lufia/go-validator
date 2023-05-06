package validator_test

import (
	"fmt"
	"io"

	"github.com/lufia/go-validator"
)

func Example_customMessage() {
	type Data struct {
		Num  int
		Name string
	}
	nameErrPrinter := validator.LengthViolationPrinterFunc[string](func(w io.Writer, min, max int) {
		fmt.Fprintf(w, "must be of length %d to %d", min, max)
	})
	v := validator.Struct(func(s validator.StructRuleAdder, r *Data) {
		s.Add(validator.Field(&r.Name, "name"),
			validator.Length[string](3, 100).WithPrinter(nameErrPrinter),
		)
	})
	err := v.Validate(&Data{
		Num:  10,
		Name: "xx",
	})
	fmt.Println(err)
	// Output:
	// name: must be of length 3 to 100
}
