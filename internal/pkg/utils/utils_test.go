package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes uint64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatBytes(tt.bytes))
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{100 * time.Microsecond, "100.00µs"},
		{100 * time.Millisecond, "100.00ms"},
		{2 * time.Second, "2.00s"},
		{90 * time.Second, "1m30s"},
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
