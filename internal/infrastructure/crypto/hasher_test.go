package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHasher_SHA256(t *testing.T) {
	h := NewDefaultHasher()
	hash1 := h.SHA256([]byte("test"))
	hash2 := h.SHA256([]byte("test"))
	hash3 := h.SHA256([]byte("different"))

	assert.Equal(t, 64, len(hash1))
	assert.Equal(t, hash1, hash2)
	assert.NotEqual(t, hash1, hash3)
}

func TestDefaultHasher_SHA512(t *testing.T) {
	h := NewDefaultHasher()
	hash := h.SHA512([]byte("test"))
	assert.Equal(t, 128, len(hash))
}

func TestDefaultHasher_BCrypt(t *testing.T) {
	h := NewDefaultHasher()
	hash, err := h.BCrypt("password", 10)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = h.CompareBCrypt(hash, "password")
	assert.NoError(t, err)

	err = h.CompareBCrypt(hash, "wrong")
	assert.Error(t, err)
}

func TestDefaultHasher_BCryptDefaultCost(t *testing.T) {
	h := NewDefaultHasher()
	hash, err := h.BCrypt("password", 3) // Below min cost
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}
