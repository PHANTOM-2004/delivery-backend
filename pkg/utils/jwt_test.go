package utils

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "190514"

func TestGenerateToken(t *testing.T) {
	account := "2252707"
	password := "woshinidie"
	claims := jwt.MapClaims{
		"account":  account,
		"password": password,
	}
	tks, err := GenerateToken(claims, secret)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(tks)
	}

	res, err := ParseToken(tks, secret)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(res["account"])
		fmt.Println(res["password"])
		if res["account"] != account || res["password"] != password {
			t.Fail()
		}
	}
}
