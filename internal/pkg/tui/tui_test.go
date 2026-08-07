package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintFunctions(_ *testing.T) {
	// These print to stdout; we just verify they don't panic
	PrintSuccess("test success")
	PrintError("test error")
	PrintWarning("test warning")
	PrintInfo("test info")
	PrintTitle("test title")
	PrintBox("test box")
}

func TestProgressBar(t *testing.T) {
	assert.Contains(t, ProgressBar(50, 100, 20), "50.0%")
	assert.Equal(t, "░░░░░░░░░░ 0.0%", ProgressBar(0, 0, 10))
	assert.Equal(t, "██████████ 100.0%", ProgressBar(100, 100, 10))
	assert.Equal(t, "██████████ 100.0%", ProgressBar(150, 100, 10))
}

func TestTable(t *testing.T) {
	headers := []string{"Name", "Value"}
	rows := [][]string{{"a", "1"}, {"b", "2"}}
	result := Table(headers, rows)
	assert.Contains(t, result, "Name")
	assert.Contains(t, result, "Value")
}

func TestTableEmpty(t *testing.T) {
	result := Table([]string{"A"}, [][]string{})
	assert.Equal(t, "No data", result)
}
