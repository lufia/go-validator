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
	inErrorFormat       = newFormat("must be a valid value in %[1]v", ByName("validValues"))
	patternErrorFormat  = newFormat("must match the pattern /%[1]v/", ByName("pattern"))
	customErrorFormat   = newFormat("must be a valid value")

	minLengthErrorFormat = newFormat("the length must be no less than %[1]d", ByName("min"))
	maxLengthErrorFormat = newFormat("the length must be no greater than %[1]d", ByName("max"))
	lengthErrorFormat    = newFormat("the length must be in range(%[1]d ... %[2]d)", ByName("min"), ByName("max"))

	minErrorFormat     = newFormat("must be no less than %[1]v", ByName("min"))
	maxErrorFormat     = newFormat("must be no greater than %[1]v", ByName("max"))
	inRangeErrorFormat = newFormat("must be in range(%[1]v ... %[2]v)", ByName("min"), ByName("max"))

	structFieldErrorFormat = newFormat("%[1]s: %[2]v", ByName("name"), ByName("error"))
)

func newFormat(key string, a ...Arg) *errorFormat {
	return &errorFormat{
		ID:   key,
		Key:  key,
		Args: a,
	}
}
