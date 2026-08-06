// Package crypto provides cryptographic operations.
package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DefaultHasher implements the Hasher interface.
type DefaultHasher struct{}

// NewDefaultHasher creates a new hasher.
func NewDefaultHasher() *DefaultHasher {
	return &DefaultHasher{}
}

// SHA256 computes SHA-256 hash.
func (h *DefaultHasher) SHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// SHA512 computes SHA-512 hash.
func (h *DefaultHasher) SHA512(data []byte) string {
	sum := sha512.Sum512(data)
	return hex.EncodeToString(sum[:])
}

// BCrypt hashes a password using bcrypt.
func (h *DefaultHasher) BCrypt(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt failed: %w", err)
	}
	return string(hash), nil
}

// CompareBCrypt compares a password with a bcrypt hash.
func (h *DefaultHasher) CompareBCrypt(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
