package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomToken generates a random string that can be used as secret tokens.
func RandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
