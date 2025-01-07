package utils

import (
	"encoding/base64"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/scrypt"
)

func Encrypt(s string, sa string) string {
	// salt: 8 bytes is
	// a good length.
	salt := []byte(sa)

	dk, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		klog.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(dk)
	return res
}
