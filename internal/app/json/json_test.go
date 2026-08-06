package json

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xef/xefcli/internal/core/logger"
)

func TestNewCommand(t *testing.T) {
	log := logger.Nop()
	cmd := NewCommand(log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "json", cmd.Use)
	assert.Len(t, cmd.Commands(), 3)
}

func TestFormatCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newFormatCmd(log)

	t.Run("format file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test.json")
		require.NoError(t, os.WriteFile(tmpFile, []byte(`{"a":1}`), 0644))

		cmd.SetArgs([]string{tmpFile})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("compact output", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test.json")
		require.NoError(t, os.WriteFile(tmpFile, []byte(`{"a":1}`), 0644))

		cmd.SetArgs([]string{tmpFile, "--compact"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "bad.json")
		require.NoError(t, os.WriteFile(tmpFile, []byte(`{invalid`), 0644))

		cmd.SetArgs([]string{tmpFile})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestValidateCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newValidateCmd(log)

	t.Run("valid json", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "valid.json")
		require.NoError(t, os.WriteFile(tmpFile, []byte(`{"a":1}`), 0644))

		cmd.SetArgs([]string{tmpFile})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "invalid.json")
		require.NoError(t, os.WriteFile(tmpFile, []byte(`{bad`), 0644))

		cmd.SetArgs([]string{tmpFile})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestDiffCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newDiffCmd(log)

	t.Run("identical files", func(t *testing.T) {
		dir := t.TempDir()
		f1 := filepath.Join(dir, "a.json")
		f2 := filepath.Join(dir, "b.json")
		require.NoError(t, os.WriteFile(f1, []byte(`{"a":1}`), 0644))
		require.NoError(t, os.WriteFile(f2, []byte(`{"a":1}`), 0644))

		cmd.SetArgs([]string{f1, f2})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("different files", func(t *testing.T) {
		dir := t.TempDir()
		f1 := filepath.Join(dir, "a.json")
		f2 := filepath.Join(dir, "b.json")
		require.NoError(t, os.WriteFile(f1, []byte(`{"a":1}`), 0644))
		require.NoError(t, os.WriteFile(f2, []byte(`{"a":2}`), 0644))

		cmd.SetArgs([]string{f1, f2})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("nested diff", func(t *testing.T) {
		dir := t.TempDir()
		f1 := filepath.Join(dir, "a.json")
		f2 := filepath.Join(dir, "b.json")
		require.NoError(t, os.WriteFile(f1, []byte(`{"a":{"b":1}}`), 0644))
		require.NoError(t, os.WriteFile(f2, []byte(`{"a":{"b":2}}`), 0644))

		cmd.SetArgs([]string{f1, f2})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestDiffValue(t *testing.T) {
	t.Run("arrays", func(t *testing.T) {
		diffs := diffValue("", []interface{}{1, 2}, []interface{}{1, 3})
		assert.NotEmpty(t, diffs)
	})

	t.Run("type mismatch", func(t *testing.T) {
		diffs := diffValue("", map[string]interface{}{}, []interface{}{})
		assert.NotEmpty(t, diffs)
	})
}
