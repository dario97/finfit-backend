package fieldvalidation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"strings"
)

type Translation struct {
	tag              string
	formattedMessage string
	override         bool
	customRegisFunc  validator.RegisterTranslationsFunc
	customTransFunc  validator.TranslationFunc
}

func registerTranslations(validate *validator.Validate, translator ut.Translator, customTranslations []Translation) error {
	_ = en2.RegisterDefaultTranslations(validate, translator)

	translations := []Translation{
		{
			tag:              LteStrDateFieldValidationTag,
			formattedMessage: "{0} must be before or equal to {1}",
			override:         false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				params := strings.Split(fe.Param(), ",")
				t, err := ut.T(fe.Tag(), fe.Field(), params[0])
				if err != nil {
					return fe.(error).Error()
				}

				return t
			},
		},
	}

	translations = append(translations, customTranslations...)

	var err error

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, translator, t.customRegisFunc, t.customTransFunc)
		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
			err = validate.RegisterTranslation(t.tag, translator, registrationFunc(t.tag, t.formattedMessage, t.override), t.customTransFunc)
		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, translator, t.customRegisFunc, translateFunc)
		} else {
			err = validate.RegisterTranslation(t.tag, translator, registrationFunc(t.tag, t.formattedMessage, t.override), translateFunc)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		return fe.(error).Error()
	}

	return t
}
