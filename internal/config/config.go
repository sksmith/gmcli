package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigPath = "config.yaml"
	templatesDir      = "templates"
	eventsDir         = "events"
)

// Load loads configuration from file.
func Load() (Config, error) {
	var cfg Config

	// Check if file exists
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist yet
		return cfg, nil
	}

	// Read and parse file
	data, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to file.
func Save(cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(defaultConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// EnsureDirectories creates necessary directories.
func EnsureDirectories() error {
	dirs := []string{templatesDir, eventsDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default template if it doesn't exist
	templatePath := filepath.Join(templatesDir, "event.md.tmpl")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		defaultTemplate := `# {{.Name}}

## Details
- Calendar: {{.CalendarName}} ({{.CalendarAbbrev}})
- Date: {{.AgeAbbrev}}{{.Year}}-{{printf "%02d" .Month}}-{{printf "%02d" .Day}}
- Days Since Year 0: {{.DaysSinceZero}}

## Description
<!-- Add event description here -->

`
		if err := os.WriteFile(templatePath, []byte(defaultTemplate), 0644); err != nil {
			return fmt.Errorf("failed to create default template: %w", err)
		}
	}

	return nil
}
