// Package network provides HTTP client implementations.
package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// DefaultHTTPClient is the default HTTP client with sensible timeouts.
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient creates a new HTTP client.
func NewDefaultHTTPClient(timeout time.Duration) *DefaultHTTPClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get performs an HTTP GET request.
func (c *DefaultHTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*interfaces.HTTPResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	h := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		if len(v) > 0 {
			h[k] = v[0]
		}
	}

	return &interfaces.HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    h,
		Body:       body,
	}, nil
}

// Post performs an HTTP POST request.
func (c *DefaultHTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*interfaces.HTTPResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	h := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		if len(v) > 0 {
			h[k] = v[0]
		}
	}

	return &interfaces.HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    h,
		Body:       respBody,
	}, nil
}

// Download downloads a file from URL to destination.
func (c *DefaultHTTPClient) Download(ctx context.Context, url, dest string, headers map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Create parent directories if they don't exist
	if err = os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Printf("error closing destination file: %v\n", cerr)
		}
	}()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Compile-time interface check.
var _ interfaces.HTTPClient = (*DefaultHTTPClient)(nil)
