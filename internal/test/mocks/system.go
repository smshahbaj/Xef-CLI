package mocks

import (
	"context"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
)

// MockSystemProvider is a test double for SystemInfoProvider.
//
//nolint:govet // fieldalignment: layout is intentional for test doubles
type MockSystemProvider struct {
	CPUInfoResult     *interfaces.CPUInfo
	CPUInfoError      error
	MemoryInfoResult  *interfaces.MemoryInfo
	MemoryInfoError   error
	DiskInfoResult    []interfaces.DiskInfo
	DiskInfoError     error
	NetworkInfoResult []interfaces.NetworkInterface
	NetworkInfoError  error
}

// CPUInfo returns mock CPU info.
func (m *MockSystemProvider) CPUInfo(_ context.Context) (*interfaces.CPUInfo, error) {
	if m.CPUInfoError != nil {
		return nil, m.CPUInfoError
	}
	return m.CPUInfoResult, nil
}

// MemoryInfo returns mock memory info.
func (m *MockSystemProvider) MemoryInfo(_ context.Context) (*interfaces.MemoryInfo, error) {
	if m.MemoryInfoError != nil {
		return nil, m.MemoryInfoError
	}
	return m.MemoryInfoResult, nil
}

// DiskInfo returns mock disk info.
func (m *MockSystemProvider) DiskInfo(_ context.Context) ([]interfaces.DiskInfo, error) {
	if m.DiskInfoError != nil {
		return nil, m.DiskInfoError
	}
	return m.DiskInfoResult, nil
}

// NetworkInfo returns mock network info.
func (m *MockSystemProvider) NetworkInfo(_ context.Context) ([]interfaces.NetworkInterface, error) {
	if m.NetworkInfoError != nil {
		return nil, m.NetworkInfoError
	}
	return m.NetworkInfoResult, nil
}

// Compile-time check.
var _ interfaces.SystemInfoProvider = (*MockSystemProvider)(nil)
