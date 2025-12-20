package jwt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AyKrimino/ScribeHost/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("mystrongsecretkey")

func CreateToken(userID uint, email string, role entity.RoleType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iss":   "ScribeHost",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAndExtractClaims(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

func GetTokenExpiry(token *jwt.Token) (time.Time, time.Duration, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, 0, fmt.Errorf("failed to assert claims to jwt.MapClaims")
	}

	expClaim, ok := claims["exp"]
	if !ok {
		return time.Time{}, 0, fmt.Errorf("'exp' claim not found in token")
	}

	var expiryUnix int64
	switch exp := expClaim.(type) {
	case float64:
		expiryUnix = int64(exp)
	case int64:
		expiryUnix = exp
	case json.Number:
		if i, err := exp.Int64(); err == nil {
			expiryUnix = i
		} else {
			return time.Time{}, 0, fmt.Errorf("invalid 'exp' claim format (json.Number): %w", err)
		}
	default:
		return time.Time{}, 0, fmt.Errorf("invalid 'exp' claim type: %T", expClaim)
	}

	expiryTime := time.Unix(expiryUnix, 0).UTC()
	durationUntilExpiry := time.Until(expiryTime)

	return expiryTime, durationUntilExpiry, nil
}

func GetUserIDFromClaims(claims jwt.MapClaims) (uint, error) {
	subClaim, ok := claims["sub"]
	if !ok {
		return 0, fmt.Errorf("'sub' claim not found")
	}

	var userID uint64
	switch sub := subClaim.(type) {
	case float64:
		userID = uint64(sub)
	case int64:
		userID = uint64(sub)
	case uint64:
		userID = sub
	case json.Number:
		if i, err := sub.Int64(); err == nil {
			userID = uint64(i)
		} else {
			return 0, fmt.Errorf("invalid 'sub' claim format (json.Number): %w", err)
		}
	case string:
		parsedID, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse 'sub' to uint")
		}
		userID = parsedID
	default:
		return 0, fmt.Errorf("invalid 'sub' claim type: %T", subClaim)
	}

	// Check if userID is positive and within the range of a 32-bit unsigned integer
	if userID > 0 && userID <= ^uint64(0)>>32 {
		return uint(userID), nil
	}
	return 0, fmt.Errorf("invalid user ID value: %d", userID)
}

func GetRoleFromClaims(claims jwt.MapClaims) (string, error) {
	roleClaim, ok := claims["role"]
	if !ok {
		return "", fmt.Errorf("'role' claim not found in token")
	}

	role, ok := roleClaim.(string)
	if !ok {
		return "", fmt.Errorf("invalid 'role' claim type: %T", roleClaim)
	}

	return role, nil
}

func GetEmailFromClaims(claims jwt.MapClaims) (string, error) {
	emailClaim, ok := claims["email"]
	if !ok {
		return "", fmt.Errorf("'email' claim not found in token")
	}

	email, ok := emailClaim.(string)
	if !ok {
		return "", fmt.Errorf("invalid 'email' claim type: %T", emailClaim)
	}

	return email, nil
}
