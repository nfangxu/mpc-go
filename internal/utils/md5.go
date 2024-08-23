package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(s string) string {
	m := md5.New()
	_, _ = io.WriteString(m, s)
	return fmt.Sprintf("%x", m.Sum(nil))
}
