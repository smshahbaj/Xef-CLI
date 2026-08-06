package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	assert.NotNil(t, loader)
	assert.NotNil(t, loader.v)
}

func TestLoad(t *testing.T) {
	loader := NewLoader()
	err := loader.Load("")
	// Should not error even if no config file exists
	assert.NoError(t, err)
}

func TestSetAndGet(t *testing.T) {
	loader := NewLoader()
	loader.Set("test.key", "value")
	assert.Equal(t, "value", loader.GetString("test.key"))
	assert.Equal(t, "value", loader.Get("test.key"))
}

func TestAllSettings(t *testing.T) {
	loader := NewLoader()
	loader.Set("key1", "value1")
	loader.Set("key2", 42)

	settings := loader.AllSettings()
	assert.NotEmpty(t, settings)
}
