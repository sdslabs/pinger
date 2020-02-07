// Package providers contains implementations of `oauth.Provider`.
// Default providers supported are: google, github.
package providers

import (
	"crypto/rand"
	"encoding/base64"
)

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
