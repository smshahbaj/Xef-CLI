package dev

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/infrastructure/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()
	cmd := NewCommand(fs, log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "dev", cmd.Use)
}

func TestProjectCreateCmd(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	log := logger.Nop()
	cmd := newProjectCreateCmd(fs, log)

	t.Run("create go project", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "testproject")
		cmd.SetArgs([]string{name, "--lang", "go"})
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.True(t, fs.Exists(name))
		assert.True(t, fs.Exists(filepath.Join(name, "go.mod")))
	})

	t.Run("create python project", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "testpy")
		cmd.SetArgs([]string{name, "--lang", "python"})
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.True(t, fs.Exists(name))
		assert.True(t, fs.Exists(filepath.Join(name, "pyproject.toml")))
	})

	t.Run("directory exists", func(t *testing.T) {
		dir := t.TempDir()
		cmd.SetArgs([]string{dir, "--lang", "go"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("unsupported language", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "test")
		cmd.SetArgs([]string{name, "--lang", "rust"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestEnvCmd(t *testing.T) {
	log := logger.Nop()
	cmd := newEnvCmd(log)

	t.Run("list format", func(t *testing.T) {
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("json format", func(t *testing.T) {
		cmd.SetArgs([]string{"--format", "json"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("export format", func(t *testing.T) {
		cmd.SetArgs([]string{"--format", "export"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestToTitle(t *testing.T) {
	assert.Equal(t, "Hello", toTitle("hello"))
	assert.Equal(t, "H", toTitle("h"))
	assert.Equal(t, "", toTitle(""))
	assert.Equal(t, "HelloWorld", toTitle("helloWorld"))
}

func TestScaffoldGoProducesParsableMain(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	root := t.TempDir()
	name := filepath.Join(root, "demo")
	requireNoError := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	requireNoError(scaffoldGo(fs, name))
	data, err := os.ReadFile(filepath.Join(name, "cmd", "demo", "main.go"))
	requireNoError(err)
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "main.go", data, parser.AllErrors); err != nil {
		t.Fatalf("generated Go source is not parsable: %v\n%s", err, data)
	}
}
