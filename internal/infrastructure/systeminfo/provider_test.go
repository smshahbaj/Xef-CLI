package systeminfo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGopsutilProvider(t *testing.T) {
	p := NewGopsutilProvider()
	assert.NotNil(t, p)
}

func TestCPUInfo(t *testing.T) {
	p := NewGopsutilProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := p.CPUInfo(ctx)
	if err != nil {
		t.Skipf("CPU info not available in this environment: %v", err)
	}
	require.NoError(t, err)
	assert.NotEmpty(t, info.Model)
	assert.GreaterOrEqual(t, info.Cores, 1)
}

func TestMemoryInfo(t *testing.T) {
	p := NewGopsutilProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := p.MemoryInfo(ctx)
	require.NoError(t, err)
	assert.Greater(t, info.Total, uint64(0))
}

func TestDiskInfo(t *testing.T) {
	p := NewGopsutilProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	disks, err := p.DiskInfo(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, disks)
}

func TestNetworkInfo(t *testing.T) {
	p := NewGopsutilProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ifaces, err := p.NetworkInfo(ctx)
	require.NoError(t, err)
	// May be empty in containerized environments
	_ = ifaces
}

func TestUptime(t *testing.T) {
	p := NewGopsutilProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uptime, err := p.Uptime(ctx)
	require.NoError(t, err)
	assert.Greater(t, uptime, uint64(0))
}
