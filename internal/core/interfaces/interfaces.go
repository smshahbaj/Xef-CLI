// Package interfaces defines core abstractions used across XefCLI.
package interfaces

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("not found")

// FileSystem abstracts file system operations.
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm uint32) error
	MkdirAll(path string, perm uint32) error
	Remove(path string) error
	RemoveAll(path string) error
	Exists(path string) bool
	IsDir(path string) bool
	ListDir(path string) ([]FileInfo, error)
	WalkDir(root string, fn func(path string, info FileInfo, err error) error) error
}

// FileInfo holds metadata about a file.
type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	IsDir   bool
	ModTime int64
	Mode    uint32
}

// Hasher provides cryptographic hashing.
type Hasher interface {
	SHA256(data []byte) string
	SHA512(data []byte) string
	BCrypt(password string, cost int) (string, error)
	CompareBCrypt(hashedPassword, password string) error
}

// HTTPClient performs HTTP operations.
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*HTTPResponse, error)
	Download(ctx context.Context, url, dest string, headers map[string]string) error
}

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// SystemInfoProvider retrieves system information.
type SystemInfoProvider interface {
	CPUInfo(ctx context.Context) (*CPUInfo, error)
	MemoryInfo(ctx context.Context) (*MemoryInfo, error)
	DiskInfo(ctx context.Context) ([]DiskInfo, error)
	NetworkInfo(ctx context.Context) ([]NetworkInterface, error)
}

// CPUInfo holds CPU information.
type CPUInfo struct {
	Model        string
	Cores        int
	Threads      int
	UsagePercent float64
	Mhz          float64
}

// MemoryInfo holds memory information.
type MemoryInfo struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

// DiskInfo holds disk information.
type DiskInfo struct {
	Path        string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
	FSType      string
}

// NetworkInterface holds network interface information.
type NetworkInterface struct {
	Name      string
	Addrs     []string
	BytesSent uint64
	BytesRecv uint64
	IsUp      bool
}
