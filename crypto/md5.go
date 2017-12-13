package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

//生成字符串的MD5值
func MD5(str ...string) string {
	h := md5.New()
	for _, s := range str {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))
}
