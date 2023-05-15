package validator_test

import (
	"context"
	"fmt"

	"github.com/lufia/go-validator"
)

func Example_customMessage() {
	type Data struct {
		Num  int
		Name string
	}
	v := validator.Struct(func(s validator.StructFieldAdder, r *Data) {
		s.Add(validator.Field(&r.Name, "name",
			validator.Length[string](3, 100).WithReferenceKey("must be of length %d to %d", validator.ByName("min"), validator.ByName("max")),
		))
	})
	ctx := validator.WithPrinter(context.Background(), nil)
	err := v.Validate(ctx, Data{
		Num:  10,
		Name: "xx",
	})
	fmt.Println(err)
	// Output:
	// name: must be of length 3 to 100
}
