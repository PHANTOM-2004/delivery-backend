package app

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func RegisterValidation(data any, rules map[string]string) {
	validate.RegisterStructValidationMapRules(rules, data)
}

func ValidateStruct(data any) error {
	//NOTE: 如果没有注册，其实并不会报错
	err := validate.Struct(data)
	return err
}
