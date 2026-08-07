package mocks

import (
	"context"
	"io"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// MockHTTPClient is a test double for HTTPClient.
type MockHTTPClient struct {
	GetResponse   *interfaces.HTTPResponse
	GetError      error
	PostResponse  *interfaces.HTTPResponse
	PostError     error
	DownloadError error
	LastURL       string
	LastDest      string
}

// Get performs a mock GET request.
func (m *MockHTTPClient) Get(_ context.Context, url string, _ map[string]string) (*interfaces.HTTPResponse, error) {
	m.LastURL = url
	if m.GetError != nil {
		return nil, m.GetError
	}
	return m.GetResponse, nil
}

// Post performs a mock POST request.
func (m *MockHTTPClient) Post(_ context.Context, url string, _ io.Reader, _ map[string]string) (*interfaces.HTTPResponse, error) {
	m.LastURL = url
	if m.PostError != nil {
		return nil, m.PostError
	}
	return m.PostResponse, nil
}

// Download performs a mock download.
func (m *MockHTTPClient) Download(_ context.Context, url, dest string, _ map[string]string) error {
	m.LastURL = url
	m.LastDest = dest
	if m.DownloadError != nil {
		return m.DownloadError
	}
	return nil
}

// Compile-time check.
var _ interfaces.HTTPClient = (*MockHTTPClient)(nil)
