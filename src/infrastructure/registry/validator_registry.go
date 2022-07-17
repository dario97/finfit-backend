package registry

import (
	"finfit-backend/src/infrastructure/interfaces/controller/validators"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
)

func GetGenericFieldsValidator() validators.FieldsValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, translator)

	return validators.NewGenericFieldsValidator(validate, translator)
}
