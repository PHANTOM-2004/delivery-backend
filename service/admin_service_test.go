package service

import (
	"fmt"
	"testing"
)

var (
	account  = "test"
	password = "123456"
)

func TestGetAdminAccessToken(t *testing.T) {
	tks := GetAdminAccessToken(account)
	fmt.Println(tks)
}
