package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := New("dev")
	assert.NotNil(t, app)
	assert.NotNil(t, app.RootCmd)
	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.Logger)
	assert.Equal(t, "xef", app.RootCmd.Use)
	assert.NotEmpty(t, app.RootCmd.Commands())
}

func TestVersionInjection(t *testing.T) {
	app := New("1.2.3")
	assert.Equal(t, "1.2.3", app.RootCmd.Version)
}

func TestExecute(t *testing.T) {
	app := New("dev")
	app.RootCmd.SetArgs([]string{"--version"})
	err := app.Execute()
	assert.NoError(t, err)
}
