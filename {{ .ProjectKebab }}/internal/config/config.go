package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"{{ .Scaffold.gomod }}/internal/paths"
)

// Config holds the application configuration loaded from a YAML file.
type Config struct {
	LogLevel string `yaml:"log_level"`
{{- if .Scaffold.feature_file_logging }}
	LogFile  string `yaml:"log_file"`
{{- end }}
}

// Default returns a Config with default values.
func Default() Config {
	return Config{
		LogLevel: "info",
	}
}

// Read loads config from the default XDG config path.
// Returns default config if the file does not exist.
func Read() (Config, error) {
	return ReadFrom(filepath.Join(paths.ConfigDir(), "config.yaml"))
}

// ReadFrom loads config from the given file path.
// Returns default config if the file does not exist.
func ReadFrom(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}

		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
