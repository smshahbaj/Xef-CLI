// Package config handles application configuration with support for
// multiple formats and priority loading.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader handles configuration loading from multiple sources.
type Loader struct {
	v *viper.Viper
}

// NewLoader creates a new configuration loader.
func NewLoader() *Loader {
	v := viper.New()
	v.SetEnvPrefix("XEF")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	return &Loader{v: v}
}

// Load reads configuration from file and environment.
func (l *Loader) Load(cfgFile string) error {
	if cfgFile != "" {
		l.v.SetConfigFile(cfgFile)
	} else {
		l.v.SetConfigName("config")
		l.v.SetConfigType("yaml")
		l.v.AddConfigPath(".")
		l.v.AddConfigPath("$HOME/.xef")
		l.v.AddConfigPath("/etc/xef/")
	}

	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}
	return nil
}

// Get returns a configuration value.
func (l *Loader) Get(key string) interface{} {
	return l.v.Get(key)
}

// GetString returns a string configuration value.
func (l *Loader) GetString(key string) string {
	return l.v.GetString(key)
}

// GetBool returns a bool configuration value.
func (l *Loader) GetBool(key string) bool {
	return l.v.GetBool(key)
}

// GetInt returns an int configuration value.
func (l *Loader) GetInt(key string) int {
	return l.v.GetInt(key)
}

// Set sets a configuration value.
func (l *Loader) Set(key string, value interface{}) {
	l.v.Set(key, value)
}

// AllSettings returns all configuration settings.
func (l *Loader) AllSettings() map[string]interface{} {
	return l.v.AllSettings()
}

// ConfigFileUsed returns the path of the config file used.
func (l *Loader) ConfigFileUsed() string {
	return l.v.ConfigFileUsed()
}
