package crypto

import (
    "encoding/base64"
)

func Base64Encode_1(str string) string {
    b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
    return b64.EncodeToString([]byte(str))
}

// base64 encode
// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
func Base64Encode(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64 decode
func Base64Decode(str string) (string, error) {
    if s, err := base64.StdEncoding.DecodeString(str);err != nil {
        return "", err
    }else{
        return string(s), nil
    }
}
