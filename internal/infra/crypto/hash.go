package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashOTP(otp string) string {
	h := sha256.New()

	h.Write([]byte(otp))
	b := h.Sum(nil)

	return hex.EncodeToString(b)
}
