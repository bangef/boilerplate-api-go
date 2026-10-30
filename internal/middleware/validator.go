package middleware

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func FormatValidationError(err error) []string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "number":
				errors = append(errors, fmt.Sprintf("Field '%s' wajib berupa angka", e.Field()))
			case "alphanum":
				errors = append(errors, fmt.Sprintf("Field '%s' wajib berupa string", e.Field()))
			case "required":
				errors = append(errors, fmt.Sprintf("Field '%s' wajib diisi", e.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("Field '%s' minimal %s karakter", e.Field(), e.Param()))
			case "gt":
				errors = append(errors, fmt.Sprintf("Field '%s' harus lebih besar dari %s", e.Field(), e.Param()))
			case "gte":
				errors = append(errors, fmt.Sprintf("Field '%s' harus lebih besar atau sama dengan %s", e.Field(), e.Param()))
			default:
				errors = append(errors, fmt.Sprintf("Field '%s' tidak memenuhi validasi '%s'", e.Field(), e.Tag()))
			}
		}
	}
	return errors
}
