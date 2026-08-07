//nolint:revive // ignore meaningless package name warning for utils
package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		want  string
		bytes uint64
	}{
		{want: "0 B", bytes: 0},
		{want: "512 B", bytes: 512},
		{want: "1.0 KB", bytes: 1024},
		{want: "1.5 KB", bytes: 1536},
		{want: "1.0 MB", bytes: 1024 * 1024},
		{want: "1.0 GB", bytes: 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatBytes(tt.bytes))
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		want string
		d    time.Duration
	}{
		{want: "100.00µs", d: 100 * time.Microsecond},
		{want: "100.00ms", d: 100 * time.Millisecond},
		{want: "2.00s", d: 2 * time.Second},
		{want: "1m30s", d: 90 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatDuration(tt.d))
		})
	}
}

func TestSafeFileName(t *testing.T) {
	assert.Equal(t, "test-file", SafeFileName("test/file"))
	assert.Equal(t, "test-file", SafeFileName("test\\file"))
	assert.Equal(t, "test-file", SafeFileName("test:file"))
}

func TestTruncateString(t *testing.T) {
	assert.Equal(t, "hello", TruncateString("hello", 10))
	assert.Equal(t, "hel...", TruncateString("hello world", 6))
}

func TestProgressPercentage(t *testing.T) {
	assert.Equal(t, 50.0, ProgressPercentage(50, 100))
	assert.Equal(t, 0.0, ProgressPercentage(50, 0))
	assert.Equal(t, 100.0, ProgressPercentage(150, 100))
}

func TestContainsString(t *testing.T) {
	assert.True(t, ContainsString([]string{"a", "b", "c"}, "b"))
	assert.False(t, ContainsString([]string{"a", "b", "c"}, "d"))
}
