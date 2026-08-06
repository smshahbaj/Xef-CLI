package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/infrastructure/filesystem"
)

func TestNewCommand(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()
	cmd := NewCommand(fs, log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "file", cmd.Use)
	assert.Len(t, cmd.Commands(), 4)
}

func TestOrganizeCmd(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()
	cmd := newOrganizeCmd(fs, log)

	t.Run("missing directory", func(t *testing.T) {
		cmd.SetArgs([]string{"/nonexistent"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("not a directory", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "file.txt")
		require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
		cmd.SetArgs([]string{tmpFile})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("organize by extension dry run", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "b.go"), []byte("b"), 0644))

		cmd.SetArgs([]string{dir, "--by", "extension", "--dry-run"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestStatsCmd(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()
	cmd := newStatsCmd(fs, log)

	t.Run("missing path", func(t *testing.T) {
		cmd.SetArgs([]string{"/nonexistent"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("valid directory", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0755))

		cmd.SetArgs([]string{dir})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestCleanCmd(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()

	t.Run("clean with dry run", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.tmp"), []byte("temp"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("keep"), 0644))

		cleanCmd := newCleanCmd(fs, log)
		cleanCmd.SetArgs([]string{dir, "--dry-run"})
		err := cleanCmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("clean removes files", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.tmp"), []byte("temp"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("keep"), 0644))

		cleanCmd := newCleanCmd(fs, log)
		cleanCmd.SetArgs([]string{dir})
		err := cleanCmd.Execute()
		assert.NoError(t, err)
		assert.False(t, fs.Exists(filepath.Join(dir, "a.tmp")))
		assert.True(t, fs.Exists(filepath.Join(dir, "b.txt")))
	})
}

func TestHashFile(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello"), 0644))

	hash1, err := hashFile(fs, path)
	require.NoError(t, err)
	assert.Equal(t, 64, len(hash1))

	hash2, err := hashFile(fs, path)
	require.NoError(t, err)
	assert.Equal(t, hash1, hash2)

	// Different content = different hash
	require.NoError(t, os.WriteFile(path, []byte("world"), 0644))
	hash3, err := hashFile(fs, path)
	require.NoError(t, err)
	assert.NotEqual(t, hash1, hash3)
}
