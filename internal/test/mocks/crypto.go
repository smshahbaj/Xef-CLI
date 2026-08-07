package mocks

import "fmt"

// MockHasher is a test double for Hasher.
type MockHasher struct {
	CompareBCryptError error
	BCryptError        error
	SHA256Result       string
	SHA512Result       string
	BCryptResult       string
}

// SHA256 returns a mock hash.
func (m *MockHasher) SHA256(_ []byte) string {
	return m.SHA256Result
}

// SHA512 returns a mock hash.
func (m *MockHasher) SHA512(_ []byte) string {
	return m.SHA512Result
}

// BCrypt returns a mock hash.
func (m *MockHasher) BCrypt(_ string, _ int) (string, error) {
	if m.BCryptError != nil {
		return "", m.BCryptError
	}
	return m.BCryptResult, nil
}

// CompareBCrypt compares passwords.
func (m *MockHasher) CompareBCrypt(hashedPassword, password string) error {
	if m.CompareBCryptError != nil {
		return m.CompareBCryptError
	}
	if hashedPassword != password {
		return fmt.Errorf("password mismatch")
	}
	return nil
}
