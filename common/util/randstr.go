package util

import (
	"math/rand"
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
