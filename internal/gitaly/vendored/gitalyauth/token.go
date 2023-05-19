package gitalyauth

import (
	"crypto/hmac"
	"crypto/sha256"
)

func hmacSign(secret []byte, message string) []byte {
	mac := hmac.New(sha256.New, secret)
	// hash.Hash never returns an error.
	_, _ = mac.Write([]byte(message))

	return mac.Sum(nil)
}
