package helper

import (
	"time"

	"github.com/AyKrimino/ScribeHost/types"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mystrongsecretkey")

func CreateToken(userId uint, email string, role types.RoleType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"role":  role,
		"iss":   "ScribeHost",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
