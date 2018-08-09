package crypto

import (
	"encoding/base64"
)

// Base64Encode1 Base64Encode1
func Base64Encode1(str string) string {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	return b64.EncodeToString([]byte(str))
}

// Base64Encode base64 encode
// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 decode
func Base64Decode(str string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(s), nil
}
