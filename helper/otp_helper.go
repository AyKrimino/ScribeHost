package helper

import (
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
