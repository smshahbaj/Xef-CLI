package filesystem

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xef/xefcli/internal/core/interfaces"
)

func TestOSFileSystem(t *testing.T) {
	fsys := NewOSFileSystem()
	tmpDir := t.TempDir()

	t.Run("Write and Read", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		data := []byte("hello world")

		err := fsys.WriteFile(path, data, 0644)
		require.NoError(t, err)

		read, err := fsys.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, data, read)
	})

	t.Run("Exists", func(t *testing.T) {
		path := filepath.Join(tmpDir, "exists.txt")
		fsys.WriteFile(path, []byte("test"), 0644)
		assert.True(t, fsys.Exists(path))
		assert.False(t, fsys.Exists(filepath.Join(tmpDir, "notexists.txt")))
	})

	t.Run("IsDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "subdir")
		fsys.MkdirAll(dir, 0755)
		assert.True(t, fsys.IsDir(dir))
		assert.False(t, fsys.IsDir(filepath.Join(dir, "file.txt")))
	})

	t.Run("ListDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "listdir")
		fsys.MkdirAll(dir, 0755)
		fsys.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
		fsys.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)

		entries, err := fsys.ListDir(dir)
		require.NoError(t, err)
		assert.Len(t, entries, 2)
	})

	t.Run("Remove", func(t *testing.T) {
		path := filepath.Join(tmpDir, "remove.txt")
		fsys.WriteFile(path, []byte("test"), 0644)
		assert.True(t, fsys.Exists(path))

		err := fsys.Remove(path)
		require.NoError(t, err)
		assert.False(t, fsys.Exists(path))
	})

	t.Run("WalkDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "walkdir")
		fsys.MkdirAll(filepath.Join(dir, "sub"), 0755)
		fsys.WriteFile(filepath.Join(dir, "root.txt"), []byte("root"), 0644)
		fsys.WriteFile(filepath.Join(dir, "sub", "nested.txt"), []byte("nested"), 0644)

		var count int
		err := fsys.WalkDir(dir, func(path string, info interfaces.FileInfo, err error) error {
			count++
			return nil
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 3)
	})
}

func TestSanitizePath(t *testing.T) {
	fsys := NewOSFileSystem()

	t.Run("empty path", func(t *testing.T) {
		_, err := fsys.sanitizePath("")
		assert.Error(t, err)
	})

	t.Run("valid path", func(t *testing.T) {
		p, err := fsys.sanitizePath("/tmp/test.txt")
		assert.NoError(t, err)
		assert.NotEmpty(t, p)
	})
}
