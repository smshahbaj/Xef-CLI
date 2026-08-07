// Package mocks provides test doubles for XefCLI interfaces.
package mocks

import (
	"sync"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// MockFileSystem is a test double for FileSystem.
type MockFileSystem struct {
	Files    map[string][]byte
	Dirs     map[string]bool
	ExistsFn func(string) bool
	mu       sync.RWMutex
}

// NewMockFileSystem creates a new mock file system.
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string][]byte),
		Dirs:  make(map[string]bool),
	}
}

// ReadFile reads a file from the mock filesystem.
func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.Files[path]
	if !ok {
		return nil, interfaces.ErrNotFound
	}
	return data, nil
}

// WriteFile writes data to the mock filesystem at path.
func (m *MockFileSystem) WriteFile(path string, data []byte, _ uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Files[path] = data
	return nil
}

// MkdirAll creates a directory path in the mock filesystem.
func (m *MockFileSystem) MkdirAll(path string, _ uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Dirs[path] = true
	return nil
}

// Remove deletes a file from the mock filesystem.
func (m *MockFileSystem) Remove(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Files, path)
	return nil
}

// RemoveAll deletes all files under the given path in the mock filesystem.
func (m *MockFileSystem) RemoveAll(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for p := range m.Files {
		if len(p) >= len(path) && p[:len(path)] == path {
			delete(m.Files, p)
		}
	}
	return nil
}

// Exists reports whether the given path exists in the mock filesystem.
func (m *MockFileSystem) Exists(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.ExistsFn != nil {
		return m.ExistsFn(path)
	}
	_, ok := m.Files[path]
	if ok {
		return true
	}
	_, ok = m.Dirs[path]
	return ok
}

// IsDir reports whether the given path is a directory in the mock filesystem.
func (m *MockFileSystem) IsDir(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.Dirs[path]
	return ok
}

// ListDir returns a list of files in the mock filesystem.
func (m *MockFileSystem) ListDir(_ string) ([]interfaces.FileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]interfaces.FileInfo, 0, len(m.Files))
	for p, data := range m.Files {
		result = append(result, interfaces.FileInfo{
			Name: p,
			Path: p,
			Size: int64(len(data)),
		})
	}
	return result, nil
}

// WalkDir walks the mock filesystem and invokes fn for each file.
func (m *MockFileSystem) WalkDir(_ string, fn func(path string, info interfaces.FileInfo, err error) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for p, data := range m.Files {
		if err := fn(p, interfaces.FileInfo{
			Name: p,
			Path: p,
			Size: int64(len(data)),
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

// Compile-time check.
var _ interfaces.FileSystem = (*MockFileSystem)(nil)
