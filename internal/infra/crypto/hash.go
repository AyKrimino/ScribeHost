package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(input string) (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write([]byte(input))
	if err != nil {
		return "", err
	}

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), nil
}
