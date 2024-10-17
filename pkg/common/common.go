package common

import (
	"errors"
	"fmt"
	"reflect"
)

func StructToStr(s any) (string, error) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", errors.New("Provided value is not a struct")
	}

	res := ""
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		res += fmt.Sprintf("%s: %v\n", field.Name, value)
	}
	return res, nil
}
