package mocks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// ServerHTTPClient is a simple HTTP client test double that performs real requests to a supplied URL.
type ServerHTTPClient struct {
	Timeout time.Duration
}

// Get performs an HTTP GET request and returns the response body.
func (m *ServerHTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*interfaces.HTTPResponse, error) {
	client := &http.Client{Timeout: m.Timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &interfaces.HTTPResponse{StatusCode: resp.StatusCode, Headers: map[string]string{"Content-Type": resp.Header.Get("Content-Type")}, Body: body}, nil
}

// Post performs an HTTP POST request and returns the response body.
func (m *ServerHTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*interfaces.HTTPResponse, error) {
	client := &http.Client{Timeout: m.Timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &interfaces.HTTPResponse{StatusCode: resp.StatusCode, Headers: map[string]string{"Content-Type": resp.Header.Get("Content-Type")}, Body: payload}, nil
}

// Download performs an HTTP download to the destination file.
func (m *ServerHTTPClient) Download(ctx context.Context, url, dest string, headers map[string]string) error {
	resp, err := m.Get(ctx, url, headers)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, resp.Body, 0o644)
}

var _ interfaces.HTTPClient = (*ServerHTTPClient)(nil)
