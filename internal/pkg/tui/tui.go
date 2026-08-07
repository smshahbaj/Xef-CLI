// Package tui provides terminal UI helpers using Lipgloss.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles holds common Lipgloss styles.
var Styles = struct {
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Success     lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Info        lipgloss.Style
	TableHeader lipgloss.Style
	TableCell   lipgloss.Style
	Box         lipgloss.Style
	Code        lipgloss.Style
}{
	Title: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginTop(1).
		MarginBottom(1).
		Padding(0, 1),
	Subtitle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A0A0A0")).
		MarginBottom(1),
	Success: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true),
	Error: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true),
	Warning: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true),
	Info: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")),
	TableHeader: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1),
	TableCell: lipgloss.NewStyle().
		Padding(0, 1),
	Box: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2),
	Code: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1),
}

// PrintSuccess prints a success message.
func PrintSuccess(msg string) {
	fmt.Println(Styles.Success.Render("✓ " + msg))
}

// PrintError prints an error message.
func PrintError(msg string) {
	fmt.Println(Styles.Error.Render("✗ " + msg))
}

// PrintWarning prints a warning message.
func PrintWarning(msg string) {
	fmt.Println(Styles.Warning.Render("⚠ " + msg))
}

// PrintInfo prints an info message.
func PrintInfo(msg string) {
	fmt.Println(Styles.Info.Render("ℹ " + msg))
}

// PrintTitle prints a title.
func PrintTitle(title string) {
	fmt.Println(Styles.Title.Render(title))
}

// PrintBox prints content in a box.
func PrintBox(content string) {
	fmt.Println(Styles.Box.Render(content))
}

// ProgressBar renders a simple progress bar.
func ProgressBar(current, total int64, width int) string {
	if total <= 0 {
		return strings.Repeat("░", width) + " 0.0%"
	}
	percent := float64(current) / float64(total)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	filled := int(percent * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s %.1f%%", bar, percent*100)
}

// Table renders a simple text table.
func Table(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return "No data"
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder
	// Header
	for i, h := range headers {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(Styles.TableHeader.Render(fmt.Sprintf("%-*s", widths[i], h)))
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", sum(widths)+3*(len(headers)-1)))
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				sb.WriteString(" | ")
			}
			if i < len(widths) {
				sb.WriteString(Styles.TableCell.Render(fmt.Sprintf("%-*s", widths[i], cell)))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
