package fieldvalidation

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type Validation struct {
	tag      string
	function validator.Func
}

func NewValidation(tag string, function validator.Func) *Validation {
	return &Validation{tag: tag, function: function}
}

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

func registerValidations(validate *validator.Validate, customValidations []Validation) error {
	validations := []Validation{
		{
			tag:      LteStrDateFieldValidationTag,
			function: LteStrDateField,
		},
	}

	validations = append(validations, customValidations...)

	for _, validation := range validations {
		err := validate.RegisterValidation(validation.tag, validation.function)
		if err != nil {
			return err
		}
	}

	return nil
}
