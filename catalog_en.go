package validator

import (
	"golang.org/x/text/language"
)

func init() {
	DefaultCatalog.SetString(language.English, requiredErrorFormat.ID, "cannot be the zero value")
	DefaultCatalog.SetString(language.English, inErrorFormat.ID, "must be a valid value in %[1]v")
	DefaultCatalog.SetString(language.English, patternErrorFormat.ID, "must match the pattern /%[1]v/")
	DefaultCatalog.SetString(language.English, customErrorFormat.ID, "must be a valid value")

	DefaultCatalog.SetString(language.English, minLengthErrorFormat.ID, "the length must be no less than %[1]d")
	DefaultCatalog.SetString(language.English, maxLengthErrorFormat.ID, "the length must be no greater than %[1]d")
	DefaultCatalog.SetString(language.English, lengthErrorFormat.ID, "the length must be in range(%[1]v ... %[2]d)")

	DefaultCatalog.SetString(language.English, minErrorFormat.ID, "must be no less than %[1]v")
	DefaultCatalog.SetString(language.English, maxErrorFormat.ID, "must be no greater than %[1]v")
	DefaultCatalog.SetString(language.English, inRangeErrorFormat.ID, "must be in range(%[1]v ... %[2]v)")

	DefaultCatalog.SetString(language.English, structFieldErrorFormat.ID, "%[1]s: %[2]v")
}
