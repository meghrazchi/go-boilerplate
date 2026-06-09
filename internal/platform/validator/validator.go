package validator

import (
	"reflect"
	"strings"

	playground "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *playground.Validate
}

func New() *Validator {
	v := playground.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validate: v}
}

func (v *Validator) ValidateStruct(payload any) map[string]string {
	if payload == nil {
		return nil
	}
	if err := v.validate.Struct(payload); err != nil {
		return FormatValidationErrors(err)
	}
	return nil
}
