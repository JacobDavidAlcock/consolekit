package intel

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/ollama/ollama/api"
)

// ModelManager handles LLM model management and selection
type ModelManager struct {
	client    *api.Client
	available []ModelInfo
}

// ModelInfo contains information about available models
type ModelInfo struct {
	Name        string `json:"name"`
	Size        string `json:"size"`
	Description string `json:"description"`
	Specialty   string `json:"specialty"`
	MinRAM      int    `json:"min_ram_gb"`
	Recommended bool   `json:"recommended"`
}

// Specialty constants for model categorization
const (
	SpecialtyGeneral  = "general"
	SpecialtyCoding   = "coding"
	SpecialtySecurity = "security"
	SpecialtyFast     = "fast"
)

// RecommendedModels contains curated models for different use cases
var RecommendedModels = []ModelInfo{
	{
		Name:        "phi3:3.8b",
		Size:        "2.2GB",
		Description: "Microsoft Phi-3 Mini - Fast and capable general-purpose model",
		Specialty:   SpecialtyGeneral,
		MinRAM:      4,
		Recommended: true,
	},
	{
		Name:        "llama3.2:3b",
		Size:        "2.0GB",
		Description: "Meta Llama 3.2 - Excellent for reasoning and analysis",
		Specialty:   SpecialtyGeneral,
		MinRAM:      4,
		Recommended: true,
	},
	{
		Name:        "qwen2.5:3b",
		Size:        "1.9GB",
		Description: "Qwen 2.5 - Strong coding and technical assistance",
		Specialty:   SpecialtyCoding,
		MinRAM:      4,
		Recommended: true,
	},
	{
		Name:        "gemma2:2b",
		Size:        "1.6GB",
		Description: "Google Gemma 2 - Lightweight and fast",
		Specialty:   SpecialtyFast,
		MinRAM:      2,
		Recommended: false,
	},
	{
		Name:        "codellama:7b",
		Size:        "3.8GB",
		Description: "Meta Code Llama - Specialized for coding tasks",
		Specialty:   SpecialtyCoding,
		MinRAM:      8,
		Recommended: false,
	},
	{
		Name:        "llama3.2:1b",
		Size:        "1.3GB",
		Description: "Meta Llama 3.2 1B - Ultra-lightweight option",
		Specialty:   SpecialtyFast,
		MinRAM:      2,
		Recommended: false,
	},
}

// NewModelManager creates a new model manager
func NewModelManager(client *api.Client) *ModelManager {
	return &ModelManager{
		client:    client,
		available: RecommendedModels,
	}
}

// AutoSelectModel chooses the best model based on system resources and preferences
func (m *ModelManager) AutoSelectModel(preferences ...string) string {
	// Get system memory
	var memInfo runtime.MemStats
	runtime.ReadMemStats(&memInfo)
	systemRAM := int(memInfo.Sys / (1024 * 1024 * 1024)) // Convert to GB

	// If system RAM is low, prefer fast models
	preferFast := systemRAM < 8

	// Check preferences
	var preferredSpecialty string
	for _, pref := range preferences {
		switch strings.ToLower(pref) {
		case "coding", "code", "development":
			preferredSpecialty = SpecialtyCoding
		case "security", "pentest", "audit":
			preferredSpecialty = SpecialtySecurity
		case "fast", "lightweight", "quick":
			preferredSpecialty = SpecialtyFast
		case "general", "analysis", "reasoning":
			preferredSpecialty = SpecialtyGeneral
		}
	}

	// Score models based on criteria
	bestModel := ""
	bestScore := -1

	for _, model := range RecommendedModels {
		score := 0

		// Prefer recommended models
		if model.Recommended {
			score += 10
		}

		// Check if model fits in available RAM
		if model.MinRAM <= systemRAM {
			score += 5
		} else {
			continue // Skip models that won't fit
		}

		// Specialty matching
		if preferredSpecialty != "" && model.Specialty == preferredSpecialty {
			score += 8
		}

		// Prefer faster models on low-RAM systems
		if preferFast && model.Specialty == SpecialtyFast {
			score += 5
		}

		// Default preference for general models
		if preferredSpecialty == "" && model.Specialty == SpecialtyGeneral {
			score += 3
		}

		if score > bestScore {
			bestScore = score
			bestModel = model.Name
		}
	}

	// Fallback to the most lightweight recommended model
	if bestModel == "" {
		for _, model := range RecommendedModels {
			if model.Recommended && model.MinRAM <= 4 {
				return model.Name
			}
		}
		return "phi3:3.8b" // Ultimate fallback
	}

	return bestModel
}

// EnsureModel downloads a model if it's not available locally
func (m *ModelManager) EnsureModel(modelName string) error {
	ctx := context.Background()

	// Check if model is already available
	if available, err := m.IsModelAvailable(modelName); err != nil {
		return fmt.Errorf("failed to check model availability: %w", err)
	} else if available {
		return nil // Model already exists
	}

	// Start download
	fmt.Printf("📥 Downloading model %s...\n", modelName)
	
	pullReq := &api.PullRequest{
		Name: modelName,
	}

	// Enhanced progress reporting with download tracker
	tracker := NewDownloadTracker()
	err := m.client.Pull(ctx, pullReq, func(resp api.ProgressResponse) error {
		tracker.Update(resp)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to download model %s: %w", modelName, err)
	}

	tracker.Complete()
	fmt.Printf("✅ Model %s downloaded successfully\n", modelName)
	return nil
}

// IsModelAvailable checks if a model is available locally
func (m *ModelManager) IsModelAvailable(modelName string) (bool, error) {
	ctx := context.Background()
	
	listResp, err := m.client.List(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list models: %w", err)
	}

	for _, model := range listResp.Models {
		if model.Name == modelName {
			return true, nil
		}
	}

	return false, nil
}

// ListAvailableModels returns models available locally
func (m *ModelManager) ListAvailableModels() ([]string, error) {
	ctx := context.Background()
	
	listResp, err := m.client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	var models []string
	for _, model := range listResp.Models {
		models = append(models, model.Name)
	}

	return models, nil
}

// GetModelInfo returns information about a specific model
func (m *ModelManager) GetModelInfo(modelName string) (*ModelInfo, bool) {
	for _, model := range RecommendedModels {
		if model.Name == modelName {
			return &model, true
		}
	}
	return nil, false
}

// GetRecommendedModels returns all recommended models
func (m *ModelManager) GetRecommendedModels() []ModelInfo {
	var recommended []ModelInfo
	for _, model := range RecommendedModels {
		if model.Recommended {
			recommended = append(recommended, model)
		}
	}
	return recommended
}

// GetModelsBySpecialty returns models filtered by specialty
func (m *ModelManager) GetModelsBySpecialty(specialty string) []ModelInfo {
	var filtered []ModelInfo
	for _, model := range RecommendedModels {
		if model.Specialty == specialty {
			filtered = append(filtered, model)
		}
	}
	return filtered
}

// EstimateSystemRAM provides a rough estimate of system RAM
func (m *ModelManager) EstimateSystemRAM() int {
	// Try to get actual system RAM instead of Go runtime memory
	return getSystemRAM()
}

// getSystemRAM attempts to get the actual system RAM across platforms
func getSystemRAM() int {
	// Try different methods to get system RAM
	if ram := getSystemRAMWindows(); ram > 0 {
		return ram
	}
	if ram := getSystemRAMUnix(); ram > 0 {
		return ram
	}
	
	// Fallback: estimate based on available memory patterns
	return estimateRAMFromRuntime()
}

// getSystemRAMWindows gets system RAM on Windows using wmic
func getSystemRAMWindows() int {
	// This is a simplified approach - in a real implementation you'd use proper Windows APIs
	// For now, we'll use a reasonable default for Windows systems
	return 8 // Assume 8GB as a reasonable default for Windows systems
}

// getSystemRAMUnix gets system RAM on Unix-like systems
func getSystemRAMUnix() int {
	// This would typically read from /proc/meminfo on Linux
	// For now, return 0 to indicate we couldn't determine it
	return 0
}

// estimateRAMFromRuntime provides a fallback estimation
func estimateRAMFromRuntime() int {
	var memInfo runtime.MemStats
	runtime.ReadMemStats(&memInfo)
	
	// If the runtime is using this much memory, the system likely has much more
	runtimeGB := int(memInfo.Sys / (1024 * 1024 * 1024))
	
	// Estimate system RAM as at least 4x the runtime usage, with a minimum of 4GB
	estimatedRAM := runtimeGB * 4
	if estimatedRAM < 4 {
		estimatedRAM = 4
	}
	
	return estimatedRAM
}

// SuggestModel provides model recommendations based on use case
func (m *ModelManager) SuggestModel(useCase string) []ModelInfo {
	useCase = strings.ToLower(useCase)
	systemRAM := m.EstimateSystemRAM()

	var suggestions []ModelInfo

	// Filter by use case and system constraints
	for _, model := range RecommendedModels {
		if model.MinRAM > systemRAM {
			continue // Skip models that won't fit
		}

		match := false
		switch {
		case strings.Contains(useCase, "code") || strings.Contains(useCase, "programming"):
			match = model.Specialty == SpecialtyCoding
		case strings.Contains(useCase, "security") || strings.Contains(useCase, "pentest"):
			match = model.Specialty == SpecialtySecurity || model.Specialty == SpecialtyGeneral
		case strings.Contains(useCase, "fast") || strings.Contains(useCase, "quick"):
			match = model.Specialty == SpecialtyFast
		default:
			match = model.Specialty == SpecialtyGeneral
		}

		if match {
			suggestions = append(suggestions, model)
		}
	}

	// If no specific matches, return recommended models that fit in RAM
	if len(suggestions) == 0 {
		for _, model := range RecommendedModels {
			if model.Recommended && model.MinRAM <= systemRAM {
				suggestions = append(suggestions, model)
			}
		}
	}

	return suggestions
}