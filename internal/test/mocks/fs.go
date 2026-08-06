// Package mocks provides test doubles for XefCLI interfaces.
package mocks

import (
	"sync"

	"github.com/xef/xefcli/internal/core/interfaces"
)

// MockFileSystem is a test double for FileSystem.
type MockFileSystem struct {
	mu       sync.RWMutex
	Files    map[string][]byte
	Dirs     map[string]bool
	ExistsFn func(string) bool
}

// NewMockFileSystem creates a new mock file system.
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string][]byte),
		Dirs:  make(map[string]bool),
	}
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.Files[path]
	if !ok {
		return nil, interfaces.ErrNotFound
	}
	return data, nil
}

func (m *MockFileSystem) WriteFile(path string, data []byte, perm uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Files[path] = data
	return nil
}

func (m *MockFileSystem) MkdirAll(path string, perm uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Dirs[path] = true
	return nil
}

func (m *MockFileSystem) Remove(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Files, path)
	return nil
}

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

func (m *MockFileSystem) IsDir(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.Dirs[path]
	return ok
}

func (m *MockFileSystem) ListDir(path string) ([]interfaces.FileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []interfaces.FileInfo
	for p, data := range m.Files {
		result = append(result, interfaces.FileInfo{
			Name: p,
			Path: p,
			Size: int64(len(data)),
		})
	}
	return result, nil
}

func (m *MockFileSystem) WalkDir(root string, fn func(path string, info interfaces.FileInfo, err error) error) error {
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
