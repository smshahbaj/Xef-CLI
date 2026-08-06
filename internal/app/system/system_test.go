package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xef/xefcli/internal/core/interfaces"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/test/mocks"
)

func TestNewCommand(t *testing.T) {
	provider := &mocks.MockSystemProvider{}
	log := logger.Nop()
	cmd := NewCommand(provider, log)
	assert.NotNil(t, cmd)
	assert.Equal(t, "system", cmd.Use)
	assert.Len(t, cmd.Commands(), 3)
}

func TestCPUCmd(t *testing.T) {
	provider := &mocks.MockSystemProvider{
		CPUInfoResult: &interfaces.CPUInfo{
			Model: "Test CPU", Cores: 4, Threads: 8, UsagePercent: 25.5, Mhz: 3200,
		},
	}
	log := logger.Nop()
	cmd := newCPUCmd(provider, log)

	t.Run("cpu info", func(t *testing.T) {
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("cpu error", func(t *testing.T) {
		provider.CPUInfoError = assert.AnError
		err := cmd.Execute()
		assert.Error(t, err)
		provider.CPUInfoError = nil
	})
}

func TestMemoryCmd(t *testing.T) {
	provider := &mocks.MockSystemProvider{
		MemoryInfoResult: &interfaces.MemoryInfo{
			Total: 16 * 1024 * 1024 * 1024, Used: 8 * 1024 * 1024 * 1024, UsedPercent: 50,
		},
	}
	log := logger.Nop()
	cmd := newMemoryCmd(provider, log)

	t.Run("memory info", func(t *testing.T) {
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("memory error", func(t *testing.T) {
		provider.MemoryInfoError = assert.AnError
		err := cmd.Execute()
		assert.Error(t, err)
		provider.MemoryInfoError = nil
	})
}

func TestDiskCmd(t *testing.T) {
	provider := &mocks.MockSystemProvider{
		DiskInfoResult: []interfaces.DiskInfo{
			{Path: "/", Total: 100 * 1024 * 1024 * 1024, Free: 50 * 1024 * 1024 * 1024, UsedPercent: 50, FSType: "ext4"},
		},
	}
	log := logger.Nop()
	cmd := newDiskCmd(provider, log)

	t.Run("disk info", func(t *testing.T) {
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("disk error", func(t *testing.T) {
		provider.DiskInfoError = assert.AnError
		err := cmd.Execute()
		assert.Error(t, err)
		provider.DiskInfoError = nil
	})
}
