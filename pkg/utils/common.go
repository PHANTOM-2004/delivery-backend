package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandString(n int) string {
	r_src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(r_src)

	var sb strings.Builder

	for i := 0; i < n; i++ {
		pos := r.Intn(len(letters))
		s := string(letters[pos])
		sb.WriteString(s)
	}

	return sb.String()
}

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
