package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Level:  InfoLevel,
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			cfg: Config{
				Level:  "invalid",
				Format: "json",
			},
			wantErr: true,
		},
		{
			name: "pretty format",
			cfg: Config{
				Level:  DebugLevel,
				Format: "pretty",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := New(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, log)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		level   LogLevel
		wantErr bool
	}{
		{DebugLevel, false},
		{InfoLevel, false},
		{WarnLevel, false},
		{ErrorLevel, false},
		{FatalLevel, false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			_, err := parseLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNop(t *testing.T) {
	log := Nop()
	assert.NotNil(t, log)
	log.Info("test message", String("key", "value"))
}

func TestFields(t *testing.T) {
	assert.Equal(t, "key", String("key", "test").Key)
	assert.Equal(t, "test", String("key", "test").Value)
	assert.Equal(t, "key", Int("key", 42).Key)
	assert.Equal(t, 42, Int("key", 42).Value)
	assert.Equal(t, "key", Int64("key", 42).Key)
	assert.Equal(t, int64(42), Int64("key", 42).Value)
	assert.Equal(t, "key", Float64("key", 3.14).Key)
	assert.Equal(t, 3.14, Float64("key", 3.14).Value)
	assert.Equal(t, "key", Bool("key", true).Key)
	assert.Equal(t, true, Bool("key", true).Value)
}
