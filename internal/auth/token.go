package auth

import (
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
)

const (
	MIN = 100000
	MAX = 999999
)

func GenerateOTP() string {
	n := rand.Intn(MAX-MIN) + MIN
	return strconv.Itoa(n)
}

func HashOTP(otp string) string {
	h := sha256.New()

	h.Write([]byte(otp))
	b := h.Sum(nil)

	return hex.EncodeToString(b)
}

func GeneratePasswordResetToken() (string, error) {
	token := make([]byte, 32)
	_, err := cryptoRand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
