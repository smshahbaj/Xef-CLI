package system

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/xef/xefcli/internal/core/interfaces"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/pkg/tui"
	"github.com/xef/xefcli/internal/pkg/utils"
)

// NewCommand creates the system command group.
func NewCommand(provider interfaces.SystemInfoProvider, log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System information tools",
		Long:  "Monitor CPU, memory, disk, and network resources.",
	}

	cmd.AddCommand(newCPUCmd(provider, log))
	cmd.AddCommand(newMemoryCmd(provider, log))
	cmd.AddCommand(newDiskCmd(provider, log))
	return cmd
}

func newCPUCmd(provider interfaces.SystemInfoProvider, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "cpu",
		Short: "Show CPU information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			info, err := provider.CPUInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get CPU info: %w", err)
			}

			tui.PrintTitle("CPU Information")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.Append([]string{"Model", info.Model})
			table.Append([]string{"Cores", fmt.Sprintf("%d", info.Cores)})
			table.Append([]string{"Threads", fmt.Sprintf("%d", info.Threads)})
			table.Append([]string{"Frequency", fmt.Sprintf("%.0f MHz", info.Mhz)})
			table.Append([]string{"Usage", fmt.Sprintf("%.1f%%", info.UsagePercent)})
			table.Render()
			return nil
		},
	}
}

func newMemoryCmd(provider interfaces.SystemInfoProvider, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "memory",
		Short: "Show memory information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			info, err := provider.MemoryInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get memory info: %w", err)
			}

			tui.PrintTitle("Memory Information")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.Append([]string{"Total", utils.FormatBytes(info.Total)})
			table.Append([]string{"Used", utils.FormatBytes(info.Used)})
			table.Append([]string{"Available", utils.FormatBytes(info.Available)})
			table.Append([]string{"Used %", fmt.Sprintf("%.1f%%", info.UsedPercent)})
			table.Render()
			return nil
		},
	}
}

func newDiskCmd(provider interfaces.SystemInfoProvider, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "disk",
		Short: "Show disk usage information",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			disks, err := provider.DiskInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get disk info: %w", err)
			}

			tui.PrintTitle("Disk Information")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Mount", "Total", "Used", "Free", "Usage %", "FS Type"})
			for _, d := range disks {
				table.Append([]string{
					d.Path,
					utils.FormatBytes(d.Total),
					utils.FormatBytes(d.Used),
					utils.FormatBytes(d.Free),
					fmt.Sprintf("%.1f%%", d.UsedPercent),
					d.FSType,
				})
			}
			table.Render()
			return nil
		},
	}
}
