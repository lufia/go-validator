package validator

import (
	"golang.org/x/text/language"
)

func init() {
	DefaultCatalog.SetString(language.Japanese, requiredErrorFormat.ID, "必須です")
	DefaultCatalog.SetString(language.Japanese, inErrorFormat.ID, "%[1]vのいずれかでなければなりません")
	DefaultCatalog.SetString(language.Japanese, patternErrorFormat.ID, "%[1]vのパターンに一致しなければなりません")
	DefaultCatalog.SetString(language.Japanese, customErrorFormat.ID, "有効な値でなければなりません")

	DefaultCatalog.SetString(language.Japanese, minLengthErrorFormat.ID, "%[1]d文字以上の長さが必要です")
	DefaultCatalog.SetString(language.Japanese, maxLengthErrorFormat.ID, "%[1]d文字以内の長さに制限されています")
	DefaultCatalog.SetString(language.Japanese, lengthErrorFormat.ID, "長さは%[1]d以上%[2]d以内の制限があります")

	DefaultCatalog.SetString(language.Japanese, minErrorFormat.ID, "%[1]v以上の値が必要です")
	DefaultCatalog.SetString(language.Japanese, maxErrorFormat.ID, "%[1]v以下の値が必要です")
	DefaultCatalog.SetString(language.Japanese, inRangeErrorFormat.ID, "%[1]v以上%[2]d以下の値が必要です")

	DefaultCatalog.SetString(language.Japanese, structFieldErrorFormat.ID, "%[1]s: %[2]v")
}
