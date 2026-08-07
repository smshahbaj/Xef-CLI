package mocks

import (
	"context"
	"io"
	"sync"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// MockHTTPClient is a test double for HTTPClient.
type MockHTTPClient struct {
	GetResponse   *interfaces.HTTPResponse
	PostResponse  *interfaces.HTTPResponse
	GetError      error
	PostError     error
	DownloadError error
	LastURL       string
	LastDest      string
	mu            sync.Mutex
}

// Get performs a mock GET request.
func (m *MockHTTPClient) Get(_ context.Context, url string, _ map[string]string) (*interfaces.HTTPResponse, error) {
	m.mu.Lock()
	m.LastURL = url
	m.mu.Unlock()

	if m.GetError != nil {
		return nil, m.GetError
	}
	return m.GetResponse, nil
}

// Post performs a mock POST request.
func (m *MockHTTPClient) Post(_ context.Context, url string, _ io.Reader, _ map[string]string) (*interfaces.HTTPResponse, error) {
	m.mu.Lock()
	m.LastURL = url
	m.mu.Unlock()

	if m.PostError != nil {
		return nil, m.PostError
	}
	return m.PostResponse, nil
}

// Download performs a mock download.
func (m *MockHTTPClient) Download(_ context.Context, url, dest string, _ map[string]string) error {
	m.mu.Lock()
	m.LastURL = url
	m.LastDest = dest
	m.mu.Unlock()

	if m.DownloadError != nil {
		return m.DownloadError
	}
	return nil
}

// GetLastURL returns the last URL requested (thread-safe).
func (m *MockHTTPClient) GetLastURL() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.LastURL
}

// GetLastDest returns the last download destination (thread-safe).
func (m *MockHTTPClient) GetLastDest() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.LastDest
}

// Compile-time check.
var _ interfaces.HTTPClient = (*MockHTTPClient)(nil)
