package fieldvalidation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

const LteStrDateFieldValidationTag = "lteStrDateField"

func LteStrDateField(fieldLevel validator.FieldLevel) bool {
	tagParam := fieldLevel.Param()
	params := strings.Split(tagParam, ",")
	fieldToCompare := params[0]
	dateFormat := params[1]
	strDate := fieldLevel.Field().String()
	strDateToCompare := fieldLevel.Parent().FieldByName(fieldToCompare).String()

	date, err := time.Parse(dateFormat, strDate)
	if err != nil {
		return false
	}

	dateToCompare, err := time.Parse(dateFormat, strDateToCompare)
	if err != nil {
		return false
	}

	return date.Before(dateToCompare) || date.Equal(dateToCompare)
}

func RegisterValidations(validate *validator.Validate, translator ut.Translator) {
	err := validate.RegisterValidation(LteStrDateFieldValidationTag, LteStrDateField)
	if err != nil {
		panic(err)
	}
	RegisterTranslations(validate, translator)
}
