package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationErrors struct {
	Errors []string `json:"errors"`
}

// Validates the request payload based on the struct's validate tags.
func Validate(payload any) error {
	structFieldTagMap := getStructFieldAndTagValues(reflect.TypeOf(payload), "yaml")

	validate := validator.New()

	err := validate.Struct(payload)
	var param string
	if err != nil {
		var errs validationErrors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid input format")
		}

		for _, err := range err.(validator.ValidationErrors) {
			param = fmt.Sprintf("%s: %s", err.Tag(), err.Param())
			if err.Param() == "" {
				param = err.Tag()
			}
			errs.Errors = append(errs.Errors, fmt.Sprintf("expected %s on field %s", param, structFieldTagMap[err.Field()]))
		}
		return errors.New(strings.Join(errs.Errors, "\n"))
	}
	return nil
}

func getStructFieldAndTagValues(payload reflect.Type, tag string) map[string]string {
	if payload.Kind() == reflect.Slice {
		return getStructFieldAndTagValues(payload.Elem(), tag)
	}

	if payload.Kind() == reflect.Ptr {
		payload = payload.Elem()
	}

	if payload.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]string)
	for i := 0; i < payload.NumField(); i++ {
		field := payload.Field(i)
		result[field.Name] = field.Tag.Get(tag)
		if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice {
			for k, v := range getStructFieldAndTagValues(field.Type, tag) {
				result[k] = v
			}
		}
	}

	return result
}
