package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xef/xefcli/internal/core/interfaces"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/test/mocks"
)

func TestNewCommand(t *testing.T) {
	client := &mocks.MockHTTPClient{}
	log := logger.Nop()
	cmd := NewCommand(client, log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "http", cmd.Use)
	assert.Len(t, cmd.Commands(), 3)
}

func TestGetCmd(t *testing.T) {
	client := &mocks.MockHTTPClient{
		GetResponse: &interfaces.HTTPResponse{StatusCode: 200, Body: []byte("ok")},
	}
	log := logger.Nop()
	cmd := newGetCmd(client, log)

	t.Run("successful get", func(t *testing.T) {
		cmd.SetArgs([]string{"https://example.com"})
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", client.LastURL)
	})

	t.Run("get error", func(t *testing.T) {
		client.GetError = assert.AnError
		cmd.SetArgs([]string{"https://example.com"})
		err := cmd.Execute()
		assert.Error(t, err)
		client.GetError = nil
	})
}

func TestDownloadCmd(t *testing.T) {
	client := &mocks.MockHTTPClient{}
	log := logger.Nop()
	cmd := newDownloadCmd(client, log)

	t.Run("download", func(t *testing.T) {
		cmd.SetArgs([]string{"https://example.com/file.zip"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestBenchmarkCmd(t *testing.T) {
	client := &mocks.MockHTTPClient{
		GetResponse: &interfaces.HTTPResponse{StatusCode: 200, Body: []byte("ok")},
	}
	log := logger.Nop()
	cmd := newBenchmarkCmd(client, log)

	t.Run("benchmark", func(t *testing.T) {
		cmd.SetArgs([]string{"https://example.com", "--requests", "10", "--concurrency", "2"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func TestParseHeaders(t *testing.T) {
	h := parseHeaders([]string{"Content-Type: application/json", "Authorization: Bearer token"})
	assert.Equal(t, "application/json", h["Content-Type"])
	assert.Equal(t, "Bearer token", h["Authorization"])
	assert.Empty(t, parseHeaders([]string{"invalid"}))
}
