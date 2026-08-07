package filesystem

import (
	"path/filepath"
	"testing"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSFileSystem(t *testing.T) {
	fsys := NewOSFileSystem()
	tmpDir := t.TempDir()

	t.Run("Write and Read", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		data := []byte("hello world")

		require.NoError(t, fsys.WriteFile(path, data, 0o644))

		read, err := fsys.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, data, read)
	})

	t.Run("Exists", func(t *testing.T) {
		path := filepath.Join(tmpDir, "exists.txt")
		require.NoError(t, fsys.WriteFile(path, []byte("test"), 0o644))
		assert.True(t, fsys.Exists(path))
		assert.False(t, fsys.Exists(filepath.Join(tmpDir, "notexists.txt")))
	})

	t.Run("IsDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "subdir")
		require.NoError(t, fsys.MkdirAll(dir, 0o755))
		assert.True(t, fsys.IsDir(dir))
		assert.False(t, fsys.IsDir(filepath.Join(dir, "file.txt")))
	})

	t.Run("ListDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "listdir")
		require.NoError(t, fsys.MkdirAll(dir, 0o755))
		require.NoError(t, fsys.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
		require.NoError(t, fsys.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644))

		entries, err := fsys.ListDir(dir)
		require.NoError(t, err)
		assert.Len(t, entries, 2)
	})

	t.Run("Remove", func(t *testing.T) {
		path := filepath.Join(tmpDir, "remove.txt")
		require.NoError(t, fsys.WriteFile(path, []byte("test"), 0o644))
		assert.True(t, fsys.Exists(path))

		err := fsys.Remove(path)
		require.NoError(t, err)
		assert.False(t, fsys.Exists(path))
	})

	t.Run("WalkDir", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "walkdir")
		require.NoError(t, fsys.MkdirAll(filepath.Join(dir, "sub"), 0o755))
		require.NoError(t, fsys.WriteFile(filepath.Join(dir, "root.txt"), []byte("root"), 0o644))
		require.NoError(t, fsys.WriteFile(filepath.Join(dir, "sub", "nested.txt"), []byte("nested"), 0o644))

		var count int
		err := fsys.WalkDir(dir, func(_ string, _ interfaces.FileInfo, err error) error {
			if err != nil {
				return err
			}
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
