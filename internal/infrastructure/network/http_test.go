package network

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultHTTPClient(t *testing.T) {
	c := NewDefaultHTTPClient(0)
	assert.NotNil(t, c)
	assert.Equal(t, 30*time.Second, c.client.Timeout)

	c2 := NewDefaultHTTPClient(5 * time.Second)
	assert.Equal(t, 5*time.Second, c2.client.Timeout)
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("hello")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewDefaultHTTPClient(5 * time.Second)
	ctx := context.Background()

	resp, err := client.Get(ctx, server.URL, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "hello", string(resp.Body))
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("created")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewDefaultHTTPClient(5 * time.Second)
	ctx := context.Background()

	resp, err := client.Post(ctx, server.URL, nil, map[string]string{"Content-Type": "application/json"})
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("file content")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewDefaultHTTPClient(5 * time.Second)
	ctx := context.Background()
	dest := filepath.Join(t.TempDir(), "subdir", "download.txt")

	err := client.Download(ctx, server.URL, dest, nil)
	require.NoError(t, err)

	data, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "file content", string(data))
}

func TestDownloadBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewDefaultHTTPClient(5 * time.Second)
	ctx := context.Background()
	err := client.Download(ctx, server.URL, filepath.Join(t.TempDir(), "file.txt"), nil)
	assert.Error(t, err)
}

type failingBody struct {
	read bool
}

func (b *failingBody) Read(p []byte) (int, error) {
	if b.read {
		return 0, io.ErrUnexpectedEOF
	}
	b.read = true
	copy(p, []byte("partial"))
	return len("partial"), nil
}

func (b *failingBody) Close() error { return nil }

func TestDownloadFailureDoesNotLeavePartialDestination(t *testing.T) {
	client := NewDefaultHTTPClient(2 * time.Second)
	client.client.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &failingBody{},
			Header:     make(http.Header),
		}, nil
	})

	dest := filepath.Join(t.TempDir(), "download.txt")
	err := client.Download(context.Background(), "https://example.test/file", dest, nil)
	if err == nil {
		t.Fatal("expected interrupted download to fail")
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("destination should not exist after failed download, stat error: %v", statErr)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
