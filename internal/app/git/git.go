// Package git provides Git utility commands.
package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/spf13/cobra"
)

// NewCommand creates the git command group.
func NewCommand(log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Git utility tools",
		Long:  "Analyze and manage Git repositories.",
	}

	cmd.AddCommand(newStatsCmd(log))
	cmd.AddCommand(newBranchesCmd(log))
	return cmd
}

func newStatsCmd(_ logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "stats",
		Short:   "Show Git repository statistics",
		Example: `  xef git stats`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if !isGitRepo(ctx) {
				return fmt.Errorf("not a git repository")
			}

			commits, err := gitCount(ctx, "HEAD")
			if err != nil {
				return fmt.Errorf("failed to count commits: %w", err)
			}

			authors, err := gitAuthors(ctx)
			if err != nil {
				return fmt.Errorf("failed to get authors: %w", err)
			}

			branch, err := gitCurrentBranch(ctx)
			if err != nil {
				branch = "unknown"
			}

			tui.PrintTitle("Git Statistics")
			fmt.Printf("Current Branch: %s\n", branch)
			fmt.Printf("Total Commits:  %s\n", commits)
			fmt.Println()

			if len(authors) > 0 {
				tui.PrintTitle("Top Contributors")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Author", "Commits"})
				for _, a := range authors {
					table.Append([]string{a.name, a.count})
				}
				table.Render()
			}

			return nil
		},
	}
}

func newBranchesCmd(_ logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "branches",
		Short:   "List branches with last commit info",
		Example: `  xef git branches`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if !isGitRepo(ctx) {
				return fmt.Errorf("not a git repository")
			}

			out, err := exec.CommandContext(ctx, "git", "branch", "-a", "--format=%(refname:short)|%(committerdate:short)|%(subject)").Output()
			if err != nil {
				return fmt.Errorf("failed to list branches: %w", err)
			}

			tui.PrintTitle("Branches")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Branch", "Last Commit", "Message"})
			scanner := bufio.NewScanner(bytes.NewReader(out))
			for scanner.Scan() {
				parts := strings.SplitN(scanner.Text(), "|", 3)
				if len(parts) == 3 {
					table.Append(parts)
				}
			}
			table.Render()
			return nil
		},
	}
}

func isGitRepo(ctx context.Context) bool {
	_, err := exec.CommandContext(ctx, "git", "rev-parse", "--git-dir").Output()
	return err == nil
}

func gitCount(ctx context.Context, rev string) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "rev-list", "--count", rev).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitCurrentBranch(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

type authorStat struct {
	name  string
	count string
}

func gitAuthors(ctx context.Context) ([]authorStat, error) {
	out, err := exec.CommandContext(ctx, "git", "shortlog", "-sn", "HEAD").Output()
	if err != nil {
		return nil, err
	}

	var authors []authorStat
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(line, "	", 2)
		if len(parts) == 2 {
			authors = append(authors, authorStat{name: parts[1], count: strings.TrimSpace(parts[0])})
		}
	}
	return authors, nil
}
