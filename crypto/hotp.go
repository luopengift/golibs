package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"time"
)

// GoogleAuth time slice is 30 second,
// if the key is invalid, return code is zero and ignore the error.
func GoogleAuth(key string) uint32 {
	if code, _, err := TotpToken(key, 30); err != nil {
		return 0
	} else {
		return code
	}
}

// TotpToken implements TOTP Algorithm, based on HOTP.
// HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))
// when C = (T -T0) / X, T0 is Unix epoch, X is Time slice length.
// TOTP = Truncate(HMAC-SHA-1(K, (T - T0) / X))
// return params: token code<uint32>, time second remain<int64>, error<error>
func TotpToken(key string, interval int64) (uint32, int64, error) {
	k, err := base32.StdEncoding.DecodeString(key)
	if err != nil {
		return 0, 0, err
	}
	epochSeconds := time.Now().Unix()

	c := bytes(epochSeconds / interval)

	// sign using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, k)
	hmacSha1.Write(c)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	truncated := binary.BigEndian.Uint32(hashParts) & 0x7FFFFFFF

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	token := truncated % 1000000

	// compute remaining seconds before next time slice
	remain := interval - (epochSeconds % interval)
	return token, remain, nil

}

func bytes(v int64) []byte {
	var (
		res    []byte
		mask   = int64(0xFF)
		shifts = [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	)
	for i := 0; i < 8; i++ {
		res = append(res, byte((v>>shifts[i])&mask))
	}
	return res
}
