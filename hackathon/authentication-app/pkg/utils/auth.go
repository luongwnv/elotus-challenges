package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateJWTToken creates a new JWT token for the given user
func GenerateJWTToken(userID uuid.UUID, username string, jwtSecret string, jwtExpireMinutes int) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(jwtExpireMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"jti":      GenerateTokenID(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
