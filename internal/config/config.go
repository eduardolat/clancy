package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration structure for Clancy.
type Config struct {
	Agent AgentConfig `yaml:"agent"`
	Loop  LoopConfig  `yaml:"loop"`
	Input InputConfig `yaml:"input"`
}

// AgentConfig defines settings for the AI agent command.
type AgentConfig struct {
	Command string            `yaml:"command"`
	Env     map[string]string `yaml:"env"`
}

// LoopConfig defines constraints and stopping criteria for the execution loop.
type LoopConfig struct {
	MaxSteps        int           `yaml:"max_steps"`
	Timeout         string        `yaml:"timeout"`
	StopPhrase      string        `yaml:"stop_phrase"`
	StopMode        string        `yaml:"stop_mode"`
	TimeoutDuration time.Duration `yaml:"-"` // Parsed duration
}

// InputConfig defines the input prompt source.
type InputConfig struct {
	Prompt string `yaml:"prompt"`
}

// Load reads the configuration from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults if necessary
	if cfg.Loop.MaxSteps == 0 {
		cfg.Loop.MaxSteps = 10 // Default safety limit
	}
	if cfg.Loop.StopMode == "" {
		cfg.Loop.StopMode = "exact"
	}
	if cfg.Loop.Timeout == "" {
		cfg.Loop.Timeout = "30m"
	}

	// Parse timeout
	duration, err := time.ParseDuration(cfg.Loop.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout format: %w", err)
	}
	cfg.Loop.TimeoutDuration = duration

	return &cfg, nil
}

// ResolvePrompt handles the "file:" prefix logic.
// If the prompt starts with "file:", it reads the content from that path.
// Otherwise, it returns the prompt as is.
func (c *Config) ResolvePrompt() (string, error) {
	prompt := c.Input.Prompt
	if strings.HasPrefix(prompt, "file:") {
		path := strings.TrimPrefix(prompt, "file:")
		path = strings.TrimSpace(path) // Clean up potential spaces

		content, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file '%s': %w", path, err)
		}
		return string(content), nil
	}
	return prompt, nil
}
