package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xef/xefcli/internal/app"
)

func TestRootCommandConstructsAndVersionWorks(t *testing.T) {
	a := app.New("test-version")
	require.NotNil(t, a)
	require.NotNil(t, a.RootCmd)
	assert.Equal(t, "xef", a.RootCmd.Use)

	cmd := a.RootCmd
	cmd.SetArgs([]string{"--version"})
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, output.String(), "test-version")
}

func TestMainInvokesRootCommand(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	os.Stderr = w
	os.Args = []string{"xef", "--version"}

	main()

	require.NoError(t, w.Close())
	output, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Contains(t, string(output), "v1.0.1")
}
