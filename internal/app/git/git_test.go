package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitCommandsWithTemporaryRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	repoDir := createTempGitRepo(t)
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	require.NoError(t, os.Chdir(repoDir))

	t.Run("stats happy path", func(t *testing.T) {
		cmd := newStatsCmd(logger.Nop())
		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("branches happy path", func(t *testing.T) {
		cmd := newBranchesCmd(logger.Nop())
		err := cmd.Execute()
		require.NoError(t, err)
	})
}

func TestGitCommandsWithEmptyRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	repoDir, err := os.MkdirTemp("", "xefcli-git-empty-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(repoDir)
	})
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	require.NoError(t, os.Chdir(repoDir))

	err = newStatsCmd(logger.Nop()).Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count commits")
}

func TestGitCommandsWithNonGitPath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	tempDir, err := os.MkdirTemp("", "xefcli-git-nonrepo-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
		_ = os.RemoveAll(tempDir)
	})
	require.NoError(t, os.Chdir(tempDir))

	err = newStatsCmd(logger.Nop()).Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a git repository")
}

func TestGitCommandsWithMissingBinary(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	missingDir := t.TempDir()
	t.Setenv("PATH", missingDir)

	err := newBranchesCmd(logger.Nop()).Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a git repository")
}

func createTempGitRepo(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	filePath := filepath.Join(repoDir, "README.md")
	require.NoError(t, os.WriteFile(filePath, []byte("hello\n"), 0o644))

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	return repoDir
}

func TestCreateTempGitRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}

	repoDir := createTempGitRepo(t)
	_, err := os.Stat(filepath.Join(repoDir, ".git"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(repoDir, "README.md"))
	assert.NoError(t, err)
	_, err = exec.Command("git", "-C", repoDir, "rev-parse", "HEAD").CombinedOutput()
	assert.NoError(t, err, fmt.Sprintf("expected commit in %s", repoDir))
}
