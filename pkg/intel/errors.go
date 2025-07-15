package intel

import (
	"fmt"
	"strings"

	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// ErrorType defines different categories of Intel errors
type ErrorType int

const (
	ErrorTypeOllama ErrorType = iota
	ErrorTypeModel
	ErrorTypeNetwork
	ErrorTypeConfig
	ErrorTypePrompt
	ErrorTypeContext
	ErrorTypeUnknown
)

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeOllama:
		return "Ollama"
	case ErrorTypeModel:
		return "Model"
	case ErrorTypeNetwork:
		return "Network"
	case ErrorTypeConfig:
		return "Config"
	case ErrorTypePrompt:
		return "Prompt"
	case ErrorTypeContext:
		return "Context"
	default:
		return "Unknown"
	}
}

// IntelError represents a structured error with context and suggestions
type IntelError struct {
	Type        ErrorType
	Code        string
	Message     string
	Cause       error
	Suggestions []string
	Context     map[string]interface{}
}

// Error implements the error interface
func (ie *IntelError) Error() string {
	return fmt.Sprintf("[%s:%s] %s", ie.Type, ie.Code, ie.Message)
}

// Display shows a user-friendly error message with suggestions
func (ie *IntelError) Display() {
	// Show the main error
	fmt.Printf("\n%s❌ %s Error:%s %s\n", 
		output.RedColor, ie.Type, output.Reset, ie.Message)
	
	// Show suggestions if available
	if len(ie.Suggestions) > 0 {
		fmt.Printf("\n%s💡 Suggestions:%s\n", output.YellowColor, output.Reset)
		for _, suggestion := range ie.Suggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
	}
	
	// Show context if available
	if len(ie.Context) > 0 {
		fmt.Printf("\n%s🔍 Context:%s\n", output.CyanColor, output.Reset)
		for key, value := range ie.Context {
			fmt.Printf("  • %s: %v\n", key, value)
		}
	}
	
	// Show underlying cause if available
	if ie.Cause != nil {
		fmt.Printf("\n%s🔧 Technical Details:%s %s\n", 
			output.CyanColor, output.Reset, ie.Cause.Error())
	}
}

// Unwrap returns the underlying error
func (ie *IntelError) Unwrap() error {
	return ie.Cause
}

// NewIntelError creates a new structured Intel error
func NewIntelError(errorType ErrorType, code, message string, cause error) *IntelError {
	return &IntelError{
		Type:        errorType,
		Code:        code,
		Message:     message,
		Cause:       cause,
		Suggestions: []string{},
		Context:     make(map[string]interface{}),
	}
}

// WithSuggestions adds suggestions to the error
func (ie *IntelError) WithSuggestions(suggestions ...string) *IntelError {
	ie.Suggestions = append(ie.Suggestions, suggestions...)
	return ie
}

// WithContext adds context information to the error
func (ie *IntelError) WithContext(key string, value interface{}) *IntelError {
	ie.Context[key] = value
	return ie
}

// Common error constructors

// NewOllamaError creates an Ollama-related error
func NewOllamaError(code, message string, cause error) *IntelError {
	err := NewIntelError(ErrorTypeOllama, code, message, cause)
	
	switch code {
	case "not_installed":
		err.WithSuggestions(
			"Install Ollama from https://ollama.com/download",
			"Run 'intel status' to check installation",
			"Try running 'intel start' to auto-install",
		)
	case "not_running":
		err.WithSuggestions(
			"Start Ollama service: 'ollama serve'",
			"Check if port 11434 is available",
			"Try 'intel start' to auto-start service",
		)
	case "connection_failed":
		err.WithSuggestions(
			"Check if Ollama is running: 'ollama list'",
			"Verify URL in config (default: http://localhost:11434)",
			"Check firewall settings",
		)
	}
	
	return err
}

// NewModelError creates a model-related error
func NewModelError(code, message string, cause error) *IntelError {
	err := NewIntelError(ErrorTypeModel, code, message, cause)
	
	switch code {
	case "not_found":
		err.WithSuggestions(
			"Check available models: 'ollama list'",
			"Pull the model: 'ollama pull <model>'",
			"Try a different model in config",
		)
	case "download_failed":
		err.WithSuggestions(
			"Check internet connection",
			"Verify model name is correct",
			"Try again later (server might be busy)",
		)
	case "incompatible":
		err.WithSuggestions(
			"Check system requirements",
			"Try a smaller model (e.g., 'phi3:3.8b')",
			"Free up disk space",
		)
	}
	
	return err
}

// NewNetworkError creates a network-related error
func NewNetworkError(code, message string, cause error) *IntelError {
	err := NewIntelError(ErrorTypeNetwork, code, message, cause)
	
	switch code {
	case "timeout":
		err.WithSuggestions(
			"Check internet connection",
			"Try increasing timeout in config",
			"Use a faster model",
		)
	case "connection_refused":
		err.WithSuggestions(
			"Check if Ollama is running",
			"Verify the service URL",
			"Check firewall settings",
		)
	}
	
	return err
}

// NewConfigError creates a configuration-related error
func NewConfigError(code, message string, cause error) *IntelError {
	err := NewIntelError(ErrorTypeConfig, code, message, cause)
	
	switch code {
	case "invalid_model":
		err.WithSuggestions(
			"Use a valid model name (e.g., 'phi3:3.8b')",
			"Check available models: 'ollama list'",
			"See recommended models in documentation",
		)
	case "invalid_url":
		err.WithSuggestions(
			"Use format: http://localhost:11434",
			"Check if port is correct",
			"Verify protocol (http/https)",
		)
	case "invalid_timeout":
		err.WithSuggestions(
			"Use positive duration (e.g., '30s')",
			"Recommended range: 10s-300s",
			"Check format: time.Duration",
		)
	}
	
	return err
}

// HandleError processes errors and returns structured IntelError
func HandleError(err error) *IntelError {
	if err == nil {
		return nil
	}
	
	// If it's already an IntelError, return as-is
	if intelErr, ok := err.(*IntelError); ok {
		return intelErr
	}
	
	// Parse common error patterns
	errStr := err.Error()
	
	// Ollama-related errors
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "11434") {
		return NewOllamaError("connection_failed", "Cannot connect to Ollama service", err)
	}
	
	if strings.Contains(errStr, "no such file") && strings.Contains(errStr, "ollama") {
		return NewOllamaError("not_installed", "Ollama is not installed", err)
	}
	
	// Model-related errors
	if strings.Contains(errStr, "model") && strings.Contains(errStr, "not found") {
		return NewModelError("not_found", "Model not found locally", err)
	}
	
	if strings.Contains(errStr, "pull") && strings.Contains(errStr, "failed") {
		return NewModelError("download_failed", "Failed to download model", err)
	}
	
	// Network-related errors
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return NewNetworkError("timeout", "Request timed out", err)
	}
	
	if strings.Contains(errStr, "connection refused") {
		return NewNetworkError("connection_refused", "Connection refused", err)
	}
	
	// Generic error
	return NewIntelError(ErrorTypeUnknown, "generic", errStr, err)
}

// RetryableError indicates if an error can be retried
func (ie *IntelError) RetryableError() bool {
	switch ie.Type {
	case ErrorTypeNetwork:
		return ie.Code == "timeout" || ie.Code == "connection_refused"
	case ErrorTypeOllama:
		return ie.Code == "not_running"
	case ErrorTypeModel:
		return ie.Code == "download_failed"
	default:
		return false
	}
}

// GetRetryDelay returns the recommended retry delay for retryable errors
func (ie *IntelError) GetRetryDelay() int {
	switch ie.Type {
	case ErrorTypeNetwork:
		return 5 // 5 seconds
	case ErrorTypeOllama:
		return 3 // 3 seconds
	case ErrorTypeModel:
		return 10 // 10 seconds
	default:
		return 1 // 1 second
	}
}

// ShowQuickHelp displays quick help for common errors
func ShowQuickHelp() {
	fmt.Printf("\n%s🆘 Quick Help:%s\n", output.BoldColor, output.Reset)
	fmt.Printf("• %sOllama not found:%s intel status, then install from https://ollama.com\n", 
		output.YellowColor, output.Reset)
	fmt.Printf("• %sService not running:%s ollama serve, or intel start\n", 
		output.YellowColor, output.Reset)
	fmt.Printf("• %sModel not found:%s ollama list, then ollama pull <model>\n", 
		output.YellowColor, output.Reset)
	fmt.Printf("• %sConnection issues:%s Check firewall, verify URL in config\n", 
		output.YellowColor, output.Reset)
	fmt.Printf("• %sTimeout errors:%s Use smaller model, increase timeout\n", 
		output.YellowColor, output.Reset)
}