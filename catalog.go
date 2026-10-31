package validator

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

var (
	defaultLanguage = language.English
	DefaultCatalog  = catalog.NewBuilder(catalog.Fallback(defaultLanguage))
	defaultPrinter  = message.NewPrinter(defaultLanguage, message.Catalog(DefaultCatalog))
)

type errorFormat struct {
	ID string // will be set only default formats

	Key  message.Reference
	Args []Arg
}

var (
	requiredErrorFormat = newFormat("cannot be the zero value")
	inErrorFormat       = newFormat("must be a valid value in %[1]v", ByPath("validValues"))
	patternErrorFormat  = newFormat("must match the pattern /%[1]v/", ByPath("pattern"))
	customErrorFormat   = newFormat("must be a valid value")

	minLengthErrorFormat = newFormat("the length must be no less than %[1]d", ByPath("min"))
	maxLengthErrorFormat = newFormat("the length must be no greater than %[1]d", ByPath("max"))
	lengthErrorFormat    = newFormat("the length must be in range(%[1]d ... %[2]d)", ByPath("min"), ByPath("max"))

	minErrorFormat     = newFormat("must be no less than %[1]v", ByPath("min"))
	maxErrorFormat     = newFormat("must be no greater than %[1]v", ByPath("max"))
	inRangeErrorFormat = newFormat("must be in range(%[1]v ... %[2]v)", ByPath("min"), ByPath("max"))

	structFieldErrorFormat = newFormat("%[1]s: %[2]v", ByPath("name"), ByPath("error"))
)

func newFormat(key string, a ...Arg) *errorFormat {
	return &errorFormat{
		ID:   key,
		Key:  key,
		Args: a,
	}
}
