package fieldvalidation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterTranslations(validate *validator.Validate, translator ut.Translator) {
	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		return "mi mama me mima"
	}

	_ = validate.RegisterTranslation(LteStrDateFieldValidationTag, translator, regTr, transFn)
}

func regTr(translator ut.Translator) error {
	return translator.Add(LteStrDateFieldValidationTag, "hola", false)
}
