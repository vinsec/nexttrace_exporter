package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Targets []Target `yaml:"targets"`
}

// Target represents a single nexttrace target configuration
type Target struct {
	Host     string        `yaml:"host"`
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
	MaxHops  int           `yaml:"max_hops"`
}

// UnmarshalYAML implements custom unmarshaling for Target to handle duration parsing
func (t *Target) UnmarshalYAML(value *yaml.Node) error {
	type rawTarget struct {
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		Interval string `yaml:"interval"`
		MaxHops  int    `yaml:"max_hops"`
	}

	var raw rawTarget
	if err := value.Decode(&raw); err != nil {
		return err
	}

	t.Host = raw.Host
	t.Name = raw.Name
	t.MaxHops = raw.MaxHops

	// Parse interval
	if raw.Interval == "" {
		t.Interval = 5 * time.Minute // Default interval
	} else {
		duration, err := time.ParseDuration(raw.Interval)
		if err != nil {
			return fmt.Errorf("invalid interval format for target %s: %w", raw.Host, err)
		}
		t.Interval = duration
	}

	// Set default max hops if not specified
	if t.MaxHops == 0 {
		t.MaxHops = 30
	}

	// Set default name if not specified
	if t.Name == "" {
		t.Name = t.Host
	}

	return nil
}

// LoadConfig loads and parses the configuration file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Targets) == 0 {
		return fmt.Errorf("no targets defined in configuration")
	}

	targetNames := make(map[string]bool)
	for i, target := range c.Targets {
		if target.Host == "" {
			return fmt.Errorf("target %d: host is required", i)
		}

		if target.Interval < time.Second {
			return fmt.Errorf("target %s: interval must be at least 1 second", target.Host)
		}

		if target.MaxHops < 1 || target.MaxHops > 64 {
			return fmt.Errorf("target %s: max_hops must be between 1 and 64", target.Host)
		}

		// Check for duplicate names
		if targetNames[target.Name] {
			return fmt.Errorf("duplicate target name: %s", target.Name)
		}
		targetNames[target.Name] = true
	}

	return nil
}
