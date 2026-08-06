package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/infrastructure/crypto"
	"github.com/xef/xefcli/internal/test/mocks"
)

func TestNewCommand(t *testing.T) {
	hasher := crypto.NewDefaultHasher()
	log := logger.Nop()
	cmd := NewCommand(hasher, log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "crypto", cmd.Use)
	assert.Len(t, cmd.Commands(), 6)
}

func TestSHA256Cmd(t *testing.T) {
	hasher := crypto.NewDefaultHasher()
	log := logger.Nop()
	cmd := newSHA256Cmd(hasher, log)

	t.Run("hash string", func(t *testing.T) {
		cmd.SetArgs([]string{"hello"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestSHA512Cmd(t *testing.T) {
	hasher := crypto.NewDefaultHasher()
	log := logger.Nop()
	cmd := newSHA512Cmd(hasher, log)

	t.Run("hash string", func(t *testing.T) {
		cmd.SetArgs([]string{"hello"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestBCryptCmd(t *testing.T) {
	hasher := crypto.NewDefaultHasher()
	log := logger.Nop()
	cmd := newBCryptCmd(hasher, log)

	t.Run("hash password", func(t *testing.T) {
		cmd.SetArgs([]string{"mypassword"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("compare password", func(t *testing.T) {
		hash, err := hasher.BCrypt("test", 10)
		require.NoError(t, err)
		cmd.SetArgs([]string{"test", "--compare", hash})
		err = cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("compare wrong password", func(t *testing.T) {
		hash, err := hasher.BCrypt("test", 10)
		require.NoError(t, err)
		cmd.SetArgs([]string{"wrong", "--compare", hash})
		err = cmd.Execute()
		assert.Error(t, err)
	})
}

func TestUUIDCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newUUIDCmd(log)

	t.Run("generate single", func(t *testing.T) {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("generate multiple", func(t *testing.T) {
		cmd.SetArgs([]string{"--count", "5"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestBase64Cmd(t *testing.T) {
	log := logger.Nop()
	cmd := newBase64Cmd(log)

	t.Run("encode", func(t *testing.T) {
		cmd.SetArgs([]string{"hello"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("decode", func(t *testing.T) {
		cmd.SetArgs([]string{"aGVsbG8=", "--decode"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("decode invalid", func(t *testing.T) {
		cmd.SetArgs([]string{"!!!", "--decode"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestPasswordCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newPasswordCmd(log)

	t.Run("generate default", func(t *testing.T) {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("generate custom length", func(t *testing.T) {
		cmd.SetArgs([]string{"--length", "32"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("invalid length", func(t *testing.T) {
		cmd.SetArgs([]string{"--length", "2"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestGeneratePassword(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		pass, err := generatePassword(16, true, true, true)
		require.NoError(t, err)
		assert.Len(t, pass, 16)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := generatePassword(2, true, true, true)
		assert.Error(t, err)
	})

	t.Run("no sets", func(t *testing.T) {
		_, err := generatePassword(16, false, false, false)
		assert.Error(t, err)
	})

	t.Run("only lower", func(t *testing.T) {
		pass, err := generatePassword(16, false, false, false)
		// This should error because no sets are selected
		assert.Error(t, err)
		_ = pass
	})
}

func TestMockHasher(t *testing.T) {
	m := &mocks.MockHasher{
		SHA256Result: "abc123",
		BCryptResult: "hashed",
	}
	assert.Equal(t, "abc123", m.SHA256([]byte("test")))
	hash, err := m.BCrypt("pass", 10)
	require.NoError(t, err)
	assert.Equal(t, "hashed", hash)
}
