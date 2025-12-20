package auth

import (
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

