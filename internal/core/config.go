package core

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the root configuration in .kevin/config.yaml
type ProjectConfig struct {
	Project ProjectMeta   `yaml:"project"`
	Agents  []AgentConfig `yaml:"agents"`
}

// ProjectMeta holds project-level metadata
type ProjectMeta struct {
	Name string `yaml:"name"`
}

// AgentConfig defines an external tool wrapper
type AgentConfig struct {
	Name         string   `yaml:"name"`
	Executor     string   `yaml:"executor"` // "docker" or "local"
	Image        string   `yaml:"image,omitempty"`
	Command      []string `yaml:"command"`
	EnvPass      []string `yaml:"env_pass,omitempty"`
	SystemPrompt string   `yaml:"system_prompt,omitempty"`
}

// LoadConfig reads and parses the config file from the specified path
func LoadConfig(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
