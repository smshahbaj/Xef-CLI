// Package systeminfo provides system information retrieval.
package systeminfo

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/xef/xefcli/internal/core/interfaces"
)

// GopsutilProvider implements SystemInfoProvider using gopsutil.
type GopsutilProvider struct{}

// NewGopsutilProvider creates a new system info provider.
func NewGopsutilProvider() *GopsutilProvider {
	return &GopsutilProvider{}
}

// CPUInfo retrieves CPU information.
func (p *GopsutilProvider) CPUInfo(ctx context.Context) (*interfaces.CPUInfo, error) {
	info, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}
	if len(info) == 0 {
		return nil, fmt.Errorf("no CPU info available")
	}

	percent, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil || len(percent) == 0 {
		percent = []float64{0}
	}

	counts, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		counts = 0
	}

	physical, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		physical = 0
	}

	return &interfaces.CPUInfo{
		Model:        info[0].ModelName,
		Cores:        physical,
		Threads:      counts,
		UsagePercent: percent[0],
		Mhz:          info[0].Mhz,
	}, nil
}

// MemoryInfo retrieves memory information.
func (p *GopsutilProvider) MemoryInfo(ctx context.Context) (*interfaces.MemoryInfo, error) {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	return &interfaces.MemoryInfo{
		Total:       vm.Total,
		Available:   vm.Available,
		Used:        vm.Used,
		UsedPercent: vm.UsedPercent,
	}, nil
}

// DiskInfo retrieves disk information.
func (p *GopsutilProvider) DiskInfo(ctx context.Context) ([]interfaces.DiskInfo, error) {
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	result := make([]interfaces.DiskInfo, 0, len(partitions))
	for _, part := range partitions {
		usage, err := disk.UsageWithContext(ctx, part.Mountpoint)
		if err != nil {
			continue
		}
		result = append(result, interfaces.DiskInfo{
			Path:        part.Mountpoint,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
			FSType:      part.Fstype,
		})
	}
	return result, nil
}

// NetworkInfo retrieves network interface information.
func (p *GopsutilProvider) NetworkInfo(ctx context.Context) ([]interfaces.NetworkInterface, error) {
	ifaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	ioCounters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		ioCounters = nil
	}

	result := make([]interfaces.NetworkInterface, 0, len(ifaces))
	for _, iface := range ifaces {
		if iface.Name == "lo" || iface.Name == "Loopback Pseudo-Interface 1" {
			continue
		}

		addrs := make([]string, 0, len(iface.Addrs))
		for _, addr := range iface.Addrs {
			addrs = append(addrs, addr.Addr)
		}

		isUp := false
		for _, flag := range iface.Flags {
			if flag == "up" {
				isUp = true
				break
			}
		}

		ni := interfaces.NetworkInterface{
			Name:  iface.Name,
			Addrs: addrs,
			IsUp:  isUp,
		}

		for _, counter := range ioCounters {
			if counter.Name == iface.Name {
				ni.BytesSent = counter.BytesSent
				ni.BytesRecv = counter.BytesRecv
				break
			}
		}

		result = append(result, ni)
	}
	return result, nil
}

// Uptime returns system uptime in seconds.
func (p *GopsutilProvider) Uptime(ctx context.Context) (uint64, error) {
	return host.UptimeWithContext(ctx)
}

// Compile-time interface check.
var _ interfaces.SystemInfoProvider = (*GopsutilProvider)(nil)
