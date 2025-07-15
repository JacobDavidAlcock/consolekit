package intel

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ConfigValidator provides validation for Intel configuration
type ConfigValidator struct {
	knownModels map[string]ModelInfo
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	// Build known models map
	knownModels := make(map[string]ModelInfo)
	for _, model := range RecommendedModels {
		knownModels[model.Name] = model
	}
	
	return &ConfigValidator{
		knownModels: knownModels,
	}
}

// ValidateConfig validates an Intel configuration
func (cv *ConfigValidator) ValidateConfig(config *Config) error {
	if config == nil {
		return NewConfigError("nil_config", "Configuration cannot be nil", nil)
	}
	
	// Validate model
	if err := cv.ValidateModel(config.Model); err != nil {
		return err
	}
	
	// Validate URL
	if err := cv.ValidateURL(config.OllamaURL); err != nil {
		return err
	}
	
	// Validate timeout
	if err := cv.ValidateTimeout(config.Timeout); err != nil {
		return err
	}
	
	// Validate context depth
	if err := cv.ValidateContextDepth(config.ContextDepth); err != nil {
		return err
	}
	
	// Validate prompts
	if err := cv.ValidatePrompts(config.CustomPrompts); err != nil {
		return err
	}
	
	return nil
}

// ValidateModel validates the model configuration
func (cv *ConfigValidator) ValidateModel(model string) error {
	if model == "" {
		return NewConfigError("empty_model", "Model name cannot be empty", nil).
			WithSuggestions(
				"Use a valid model name (e.g., 'phi3:3.8b')",
				"Check available models: 'ollama list'",
			)
	}
	
	// Check if it's a valid model format
	if !cv.isValidModelFormat(model) {
		return NewConfigError("invalid_model_format", 
			fmt.Sprintf("Invalid model format: %s", model), nil).
			WithSuggestions(
				"Use format: 'model:version' (e.g., 'phi3:3.8b')",
				"Or just model name: 'phi3' for latest version",
			)
	}
	
	// Check if it's a known model
	if _, isKnown := cv.knownModels[model]; !isKnown {
		return NewConfigError("unknown_model", 
			fmt.Sprintf("Unknown model: %s", model), nil).
			WithSuggestions(
				"Check available models: 'ollama list'",
				"Try a recommended model: 'phi3:3.8b', 'llama3.2:3b'",
				"Visit https://ollama.com/library for more models",
			)
	}
	
	return nil
}

// ValidateURL validates the Ollama URL
func (cv *ConfigValidator) ValidateURL(ollamaURL string) error {
	if ollamaURL == "" {
		return NewConfigError("empty_url", "Ollama URL cannot be empty", nil).
			WithSuggestions(
				"Use default: 'http://localhost:11434'",
				"Or your custom Ollama server URL",
			)
	}
	
	// Parse URL
	parsedURL, err := url.Parse(ollamaURL)
	if err != nil {
		return NewConfigError("invalid_url", 
			fmt.Sprintf("Invalid URL format: %s", ollamaURL), err).
			WithSuggestions(
				"Use format: 'http://localhost:11434'",
				"Include protocol (http:// or https://)",
				"Check for typos in hostname or port",
			)
	}
	
	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return NewConfigError("invalid_url_scheme", 
			fmt.Sprintf("Invalid URL scheme: %s", parsedURL.Scheme), nil).
			WithSuggestions(
				"Use http:// or https://",
				"Example: 'http://localhost:11434'",
			)
	}
	
	// Check if port is reasonable
	if parsedURL.Port() != "" {
		if parsedURL.Port() == "0" {
			return NewConfigError("invalid_port", "Port cannot be 0", nil).
				WithSuggestions(
					"Use default port 11434",
					"Or specify a valid port number",
				)
		}
	}
	
	return nil
}

// ValidateTimeout validates the timeout configuration
func (cv *ConfigValidator) ValidateTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return NewConfigError("invalid_timeout", 
			"Timeout must be positive", nil).
			WithSuggestions(
				"Use a positive duration (e.g., '30s')",
				"Recommended range: 10s-300s",
			)
	}
	
	if timeout < 5*time.Second {
		return NewConfigError("timeout_too_short", 
			"Timeout is too short, may cause frequent failures", nil).
			WithSuggestions(
				"Use at least 5 seconds",
				"Recommended minimum: 10s",
			)
	}
	
	if timeout > 10*time.Minute {
		return NewConfigError("timeout_too_long", 
			"Timeout is too long, may cause poor user experience", nil).
			WithSuggestions(
				"Use less than 10 minutes",
				"Recommended maximum: 5m",
			)
	}
	
	return nil
}

// ValidateContextDepth validates the context depth configuration
func (cv *ConfigValidator) ValidateContextDepth(depth int) error {
	if depth < 1 {
		return NewConfigError("invalid_context_depth", 
			"Context depth must be at least 1", nil).
			WithSuggestions(
				"Use a positive number",
				"Recommended range: 5-20",
			)
	}
	
	if depth > 100 {
		return NewConfigError("context_depth_too_large", 
			"Context depth is too large, may cause memory issues", nil).
			WithSuggestions(
				"Use less than 100",
				"Recommended maximum: 50",
			)
	}
	
	return nil
}

// ValidatePrompts validates custom prompts
func (cv *ConfigValidator) ValidatePrompts(prompts map[string]string) error {
	if prompts == nil {
		return nil // OK to have no custom prompts
	}
	
	validPromptTypes := map[string]bool{
		"analyze": true,
		"suggest": true,
		"explain": true,
		"debug":   true,
		"help":    true,
	}
	
	for promptType, promptText := range prompts {
		// Check if prompt type is valid
		if !validPromptTypes[promptType] {
			return NewConfigError("invalid_prompt_type", 
				fmt.Sprintf("Invalid prompt type: %s", promptType), nil).
				WithSuggestions(
					"Valid types: analyze, suggest, explain, debug, help",
					"Check for typos in prompt type",
				)
		}
		
		// Check if prompt is not empty
		if strings.TrimSpace(promptText) == "" {
			return NewConfigError("empty_prompt", 
				fmt.Sprintf("Prompt for '%s' cannot be empty", promptType), nil).
				WithSuggestions(
					"Provide a meaningful prompt",
					"Or remove the empty prompt entry",
				)
		}
		
		// Check if prompt is reasonable length
		if len(promptText) > 1000 {
			return NewConfigError("prompt_too_long", 
				fmt.Sprintf("Prompt for '%s' is too long (%d chars)", promptType, len(promptText)), nil).
				WithSuggestions(
					"Keep prompts under 1000 characters",
					"Break long prompts into multiple lines",
				)
		}
	}
	
	return nil
}

// isValidModelFormat checks if the model name follows valid format
func (cv *ConfigValidator) isValidModelFormat(model string) bool {
	// Allow simple model names (e.g., "phi3") or versioned (e.g., "phi3:3.8b")
	validFormat := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+(:[\w\-_.]+)?$`)
	return validFormat.MatchString(model)
}

// ValidateAndNormalize validates and normalizes configuration
func (cv *ConfigValidator) ValidateAndNormalize(config *Config) error {
	if err := cv.ValidateConfig(config); err != nil {
		return err
	}
	
	// Normalize values
	config.Model = strings.ToLower(config.Model)
	config.OllamaURL = strings.TrimRight(config.OllamaURL, "/")
	
	// Set defaults for missing values
	if config.SystemPrompt == "" {
		config.SystemPrompt = "You are a concise CLI assistant. Respond like a skilled colleague - brief, direct, actionable. No fluff."
	}
	
	if config.CustomPrompts == nil {
		config.CustomPrompts = make(map[string]string)
	}
	
	return nil
}

// GetModelRecommendations returns model recommendations based on system
func (cv *ConfigValidator) GetModelRecommendations() []ModelInfo {
	recommended := make([]ModelInfo, 0)
	for _, model := range RecommendedModels {
		if model.Recommended {
			recommended = append(recommended, model)
		}
	}
	return recommended
}

// GetModelInfo returns information about a specific model
func (cv *ConfigValidator) GetModelInfo(modelName string) (*ModelInfo, bool) {
	info, exists := cv.knownModels[modelName]
	return &info, exists
}

// SuggestModel suggests a model based on preferences
func (cv *ConfigValidator) SuggestModel(preferences ...string) string {
	manager := NewModelManager(nil)
	return manager.AutoSelectModel(preferences...)
}

// ValidateSystemRequirements checks if system meets requirements for model
func (cv *ConfigValidator) ValidateSystemRequirements(modelName string) error {
	info, exists := cv.knownModels[modelName]
	if !exists {
		return NewModelError("unknown_model", 
			fmt.Sprintf("Unknown model: %s", modelName), nil)
	}
	
	// Check RAM requirements (rough estimation)
	manager := NewModelManager(nil)
	estimatedRAM := manager.EstimateSystemRAM()
	
	if estimatedRAM < info.MinRAM {
		return NewModelError("insufficient_ram", 
			fmt.Sprintf("Model requires %dGB RAM, but system has ~%dGB", 
				info.MinRAM, estimatedRAM), nil).
			WithSuggestions(
				fmt.Sprintf("Try a smaller model (e.g., 'phi3:3.8b' needs 4GB)"),
				"Free up memory by closing other applications",
				"Consider upgrading your system",
			).
			WithContext("required_ram", info.MinRAM).
			WithContext("available_ram", estimatedRAM)
	}
	
	return nil
}

// GetValidationSummary returns a summary of validation rules
func (cv *ConfigValidator) GetValidationSummary() string {
	return `Intel Configuration Validation Rules:

Model:
- Must use valid format: 'model:version' or 'model'
- Should be from recommended list
- Must meet system RAM requirements

URL:
- Must include protocol (http:// or https://)
- Should be valid URL format
- Port must be reasonable (if specified)

Timeout:
- Must be positive duration
- Recommended range: 10s-300s
- Should balance speed vs reliability

Context Depth:
- Must be positive integer
- Recommended range: 5-20
- Higher values use more memory

Custom Prompts:
- Must be valid prompt types
- Cannot be empty
- Should be under 1000 characters

Use 'intel help errors' for troubleshooting.`
}