package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// GenerateTokenID generates a unique token ID
func GenerateTokenID() string {
	return uuid.New().String()
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetCurrentTimestamp returns current unix timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
