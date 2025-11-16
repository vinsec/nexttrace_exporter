package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	content := `
targets:
  - host: 8.8.8.8
    name: test_target
    interval: 5m
    max_hops: 30
    nexttrace_args: []
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(cfg.Targets))
	}

	target := cfg.Targets[0]
	if target.Host != "8.8.8.8" {
		t.Errorf("Expected host 8.8.8.8, got %s", target.Host)
	}
	if target.Name != "test_target" {
		t.Errorf("Expected name test_target, got %s", target.Name)
	}
	if target.Interval != 5*time.Minute {
		t.Errorf("Expected interval 5m, got %v", target.Interval)
	}
	if target.MaxHops != 30 {
		t.Errorf("Expected max_hops 30, got %d", target.MaxHops)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Targets: []Target{
					{
						Host:     "8.8.8.8",
						Name:     "test",
						Interval: 5 * time.Minute,
						MaxHops:  30,
					},
				},
			},
			expectErr: false,
		},
		{
			name: "no targets",
			config: Config{
				Targets: []Target{},
			},
			expectErr: true,
		},
		{
			name: "missing host",
			config: Config{
				Targets: []Target{
					{
						Name:     "test",
						Interval: 5 * time.Minute,
						MaxHops:  30,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "interval too small",
			config: Config{
				Targets: []Target{
					{
						Host:     "8.8.8.8",
						Name:     "test",
						Interval: 500 * time.Millisecond,
						MaxHops:  30,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid max_hops",
			config: Config{
				Targets: []Target{
					{
						Host:     "8.8.8.8",
						Name:     "test",
						Interval: 5 * time.Minute,
						MaxHops:  100,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "duplicate target names",
			config: Config{
				Targets: []Target{
					{
						Host:     "8.8.8.8",
						Name:     "test",
						Interval: 5 * time.Minute,
						MaxHops:  30,
					},
					{
						Host:     "1.1.1.1",
						Name:     "test",
						Interval: 5 * time.Minute,
						MaxHops:  30,
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTargetDefaults(t *testing.T) {
	content := `
targets:
  - host: 8.8.8.8
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	target := cfg.Targets[0]

	// Check defaults
	if target.Name != "8.8.8.8" {
		t.Errorf("Expected default name to be host, got %s", target.Name)
	}
	if target.Interval != 5*time.Minute {
		t.Errorf("Expected default interval 5m, got %v", target.Interval)
	}
	if target.MaxHops != 30 {
		t.Errorf("Expected default max_hops 30, got %d", target.MaxHops)
	}
}
