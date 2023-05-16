package validator_test

import (
	"context"
	"fmt"

	"github.com/lufia/go-validator"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.English, "must be of length %[1]d to %[2]d", "must be of length %[1]d to %[2]d")
	message.SetString(language.Japanese, "must be of length %[1]d to %[2]d", "%[1]d文字以上%[2]d文字以内で入力してください")
}

func Example_customMessage() {
	type Data struct {
		Num  int
		Name string
	}
	v := validator.Struct(func(s validator.StructRule, r *Data) {
		validator.AddField(s, &r.Name, "name",
			validator.Length[string](3, 100).WithFormat("must be of length %[1]d to %[2]d", validator.ByName("min"), validator.ByName("max")),
		)
	})
	p := message.NewPrinter(language.English)
	ctx := validator.WithPrinter(context.Background(), p)
	err := v.Validate(ctx, &Data{
		Num:  10,
		Name: "xx",
	})
	fmt.Println(err)
	// Output:
	// name: must be of length 3 to 100
}
