package oauthutil

import (
	"crypto/rand"
	"encoding/base64"
)

// 32 byte の乱数を URL-safe Base64 でエンコード
func NewRandomString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
