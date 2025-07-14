package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents application configuration
type Config struct {
	data map[string]interface{}
}

// New creates a new config instance
func New() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

// LoadFromFile loads configuration from a YAML file
func (c *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return yaml.Unmarshal(data, &c.data)
}

// SaveToFile saves configuration to a YAML file
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.data[key] = value
}

// Get gets a configuration value
func (c *Config) Get(key string) (interface{}, bool) {
	value, exists := c.data[key]
	return value, exists
}

// GetString gets a string configuration value
func (c *Config) GetString(key string) (string, bool) {
	value, exists := c.data[key]
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt gets an integer configuration value
func (c *Config) GetInt(key string) (int, bool) {
	value, exists := c.data[key]
	if !exists {
		return 0, false
	}
	if i, ok := value.(int); ok {
		return i, true
	}
	return 0, false
}

// GetBool gets a boolean configuration value
func (c *Config) GetBool(key string) (bool, bool) {
	value, exists := c.data[key]
	if !exists {
		return false, false
	}
	if b, ok := value.(bool); ok {
		return b, true
	}
	return false, false
}

// LoadFromStruct loads configuration from a struct using YAML tags
func (c *Config) LoadFromStruct(v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	return yaml.Unmarshal(data, &c.data)
}

// UnmarshalTo unmarshals configuration data to a struct
func (c *Config) UnmarshalTo(v interface{}) error {
	data, err := yaml.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}

	return yaml.Unmarshal(data, v)
}

// HandleStartupFlag processes the --config startup flag
func HandleStartupFlag() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to a YAML configuration file")
	flag.Parse()

	if configPath != "" {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return "", fmt.Errorf("config file does not exist: %s", configPath)
		}
	}

	return configPath, nil
}

// GenerateExample generates an example configuration file content
func GenerateExample(appName string) string {
	return fmt.Sprintf(`# %s configuration file
#
# Use this file to pre-load your settings at startup.
# Launch with: ./%s --config /path/to/your/config.yaml

# Example configuration values
example_setting: "default_value"
debug: false
port: 8080
`, appName, appName)
}