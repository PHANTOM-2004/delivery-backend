package utils

import (
	"encoding/base64"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

func Encrypt(s string, sa string) string {
	// salt: 8 bytes is
	// a good length.
	salt := []byte(sa)

	dk, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Panic(err)
	}
	res := base64.StdEncoding.EncodeToString(dk)
	return res
}
