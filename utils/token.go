package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	// URLとして使える形式
	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}
