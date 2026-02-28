package validator

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	translationEN "github.com/go-playground/validator/v10/translations/en"
)

var (
	v *validator.Validate
	t ut.Translator
)

func Initialize() {
	v = validator.New()
	en := en.New()
	t, _ = ut.New(en, en).GetTranslator("en")
	translationEN.RegisterDefaultTranslations(v, t)
}

func ValidateStruct(e any) error {
	if v == nil {
		panic("validator is yet to initialized")
	}

	err := v.Struct(e)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			msgs := slices.Collect(maps.Values(validationErrs.Translate(t)))
			err = fmt.Errorf("validation error: %s", strings.Join(msgs, ", "))
		}
	}

	return err
}
