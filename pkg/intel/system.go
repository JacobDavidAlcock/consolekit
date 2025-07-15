package intel

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// IntelSystem is the core AI assistant system for ConsoleKit
type IntelSystem struct {
	appName        string
	client         *api.Client
	model          string
	context        *Context
	contextManager *ContextManager
	providers      []ContextProvider
	config         *Config
	ollamaManager  *OllamaManager
	initialized    bool
	mu             sync.RWMutex
}

// Config holds configuration for the Intel system
type Config struct {
	Model         string            `yaml:"model"`
	AutoDownload  bool              `yaml:"auto_download"`
	Proactive     bool              `yaml:"proactive"`
	ContextDepth  int               `yaml:"context_depth"`
	SystemPrompt  string            `yaml:"system_prompt"`
	CustomPrompts map[string]string `yaml:"custom_prompts"`
	OllamaURL     string            `yaml:"ollama_url"`
	Timeout       time.Duration     `yaml:"timeout"`
}

// Context holds the current session context for AI analysis
type Context struct {
	RecentActions []Action
	SessionData   map[string]interface{}
	StartTime     time.Time
	mu            sync.RWMutex
}

// Response represents a response from the AI system
type Response struct {
	Content   string                 `json:"content"`
	Type      string                 `json:"type"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// Suggestions represents AI-generated suggestions
type Suggestions struct {
	NextSteps   []string               `json:"next_steps"`
	Commands    []string               `json:"commands"`
	Tips        []string               `json:"tips"`
	Warnings    []string               `json:"warnings"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Explanation represents an AI explanation of concepts
type Explanation struct {
	Topic       string                 `json:"topic"`
	Summary     string                 `json:"summary"`
	Details     string                 `json:"details"`
	Examples    []string               `json:"examples"`
	References  []string               `json:"references"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Model:        "phi3:3.8b",
		AutoDownload: true,
		Proactive:    false,
		ContextDepth: 10,
		SystemPrompt: "You are a concise CLI assistant. Respond like a skilled colleague - brief, direct, actionable. No fluff.",
		CustomPrompts: map[string]string{
			"analyze": "Analyze the session. Give 3-5 key findings and immediate next steps. Use bullets. Be concise.",
			"suggest": "Suggest 3-5 specific commands to run next. Focus on actionable steps. Use bullets and code examples.",
			"explain": "Explain this concept concisely. Include key risks and 2-3 practical examples. Keep it brief.",
		},
		OllamaURL: "http://localhost:11434",
		Timeout:   30 * time.Second,
	}
}

// New creates a new Intel system
func New(appName string, config *Config) *IntelSystem {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate and normalize configuration
	validator := NewConfigValidator()
	if err := validator.ValidateAndNormalize(config); err != nil {
		// Log validation error but don't fail - use defaults
		fmt.Printf("%s⚠️  Config validation warning: %s%s\n", 
			output.YellowColor, err.Error(), output.Reset)
		config = DefaultConfig()
	}

	return &IntelSystem{
		appName:        appName,
		config:         config,
		ollamaManager:  NewOllamaManager(),
		contextManager: NewContextManager(4000), // 4000 token limit
		context:        &Context{
			RecentActions: make([]Action, 0),
			SessionData:   make(map[string]interface{}),
			StartTime:     time.Now(),
		},
		providers: make([]ContextProvider, 0),
	}
}

// Initialize sets up the Intel system and connects to Ollama
func (i *IntelSystem) Initialize() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.initialized {
		return nil
	}

	// Ensure Ollama is available (install/start if needed)
	if err := i.ollamaManager.EnsureOllamaAvailable(); err != nil {
		intelErr := HandleError(err)
		intelErr.Display()
		return intelErr
	}

	// Create Ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		intelErr := NewOllamaError("client_creation_failed", "Failed to create Ollama client", err)
		return intelErr
	}

	i.client = client

	// Validate system requirements for model
	validator := NewConfigValidator()
	if err := validator.ValidateSystemRequirements(i.config.Model); err != nil {
		if intelErr, ok := err.(*IntelError); ok {
			// Show warning but continue with model download
			fmt.Printf("%s⚠️  %s%s\n", output.YellowColor, intelErr.Message, output.Reset)
		}
	}

	// Check if model exists and download if needed
	if i.config.AutoDownload {
		if err := i.ensureModel(); err != nil {
			intelErr := HandleError(err)
			return intelErr
		}
	}

	i.initialized = true
	return nil
}

// RegisterProvider adds a context provider to the system
func (i *IntelSystem) RegisterProvider(provider ContextProvider) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.providers = append(i.providers, provider)
}

// IsInitialized returns whether the system is initialized
func (i *IntelSystem) IsInitialized() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.initialized
}

// Analyze performs AI analysis of the current session
func (i *IntelSystem) Analyze(userPrompt string) (*Response, error) {
	if !i.IsInitialized() {
		return nil, fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(userPrompt, PromptAnalyze)
	content, err := i.queryModel(prompt)
	if err != nil {
		return nil, err
	}

	return &Response{
		Content:   content,
		Type:      "analysis",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"model":        i.config.Model,
			"prompt_type":  "analyze",
			"provider_count": len(i.providers),
		},
	}, nil
}

// Suggest provides AI-generated suggestions for next steps
func (i *IntelSystem) Suggest(context string) (*Suggestions, error) {
	if !i.IsInitialized() {
		return nil, fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(context, PromptSuggest)
	content, err := i.queryModel(prompt)
	if err != nil {
		return nil, err
	}

	// Parse suggestions from the response
	suggestions := &Suggestions{
		NextSteps: []string{},
		Commands:  []string{},
		Tips:      []string{},
		Warnings:  []string{},
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"model":       i.config.Model,
			"prompt_type": "suggest",
		},
	}

	// For now, put the entire response in NextSteps
	// TODO: Parse structured output when implementing specific response formats
	suggestions.NextSteps = append(suggestions.NextSteps, content)

	return suggestions, nil
}

// Explain provides detailed explanations of concepts or findings
func (i *IntelSystem) Explain(topic string) (*Explanation, error) {
	if !i.IsInitialized() {
		return nil, fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(topic, PromptExplain)
	content, err := i.queryModel(prompt)
	if err != nil {
		return nil, err
	}

	return &Explanation{
		Topic:     topic,
		Summary:   content, // TODO: Parse structured response
		Details:   content,
		Examples:  []string{},
		References: []string{},
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"model":       i.config.Model,
			"prompt_type": "explain",
		},
	}, nil
}

// AddAction records a command action for context
func (i *IntelSystem) AddAction(command string, args []string, result string, success bool) {
	i.context.mu.Lock()
	defer i.context.mu.Unlock()

	action := Action{
		Command:   command,
		Args:      args,
		Result:    result,
		Success:   success,
		Timestamp: time.Now(),
	}

	i.context.RecentActions = append(i.context.RecentActions, action)

	// Keep only the most recent actions based on ContextDepth
	if len(i.context.RecentActions) > i.config.ContextDepth {
		i.context.RecentActions = i.context.RecentActions[len(i.context.RecentActions)-i.config.ContextDepth:]
	}
}

// ensureModel checks if the model exists and downloads it if needed
func (i *IntelSystem) ensureModel() error {
	ctx := context.Background()
	
	// List available models
	listResp, err := i.client.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	// Check if our model is available
	for _, model := range listResp.Models {
		if model.Name == i.config.Model {
			return nil // Model is available
		}
	}

	// Model not found, attempt to pull it
	ShowPersonalityMessage("downloading")
	fmt.Printf("📥 Downloading model %s...\n", i.config.Model)
	
	pullReq := &api.PullRequest{
		Name: i.config.Model,
	}

	// Enhanced progress reporting with download tracker
	tracker := NewDownloadTracker()
	err = i.client.Pull(ctx, pullReq, func(resp api.ProgressResponse) error {
		tracker.Update(resp)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}
	
	tracker.Complete()
	fmt.Printf("✅ Model %s downloaded successfully\n", i.config.Model)
	return nil
}

// queryModel sends a query to the LLM and returns the response with retry logic
func (i *IntelSystem) queryModel(prompt string) (string, error) {
	const maxRetries = 3
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), i.config.Timeout)
		
		req := &api.ChatRequest{
			Model: i.config.Model,
			Messages: []api.Message{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		}

		var response strings.Builder
		err := i.client.Chat(ctx, req, func(resp api.ChatResponse) error {
			response.WriteString(resp.Message.Content)
			return nil
		})
		
		cancel()

		if err == nil {
			return response.String(), nil
		}
		
		// Handle error with retry logic
		intelErr := HandleError(err)
		if !intelErr.RetryableError() || attempt == maxRetries {
			return "", intelErr
		}
		
		// Wait before retry
		delay := time.Duration(intelErr.GetRetryDelay()) * time.Second
		fmt.Printf("%s⏳ Retrying in %v... (attempt %d/%d)%s\n", 
			output.YellowColor, delay, attempt, maxRetries, output.Reset)
		time.Sleep(delay)
	}

	return "", NewNetworkError("max_retries", "Maximum retry attempts exceeded", nil)
}

// queryModelWithStreaming sends a query to the LLM and streams the response with formatting
func (i *IntelSystem) queryModelWithStreaming(prompt string, onToken func(string), onComplete func(string)) error {
	ctx, cancel := context.WithTimeout(context.Background(), i.config.Timeout)
	defer cancel()

	req := &api.ChatRequest{
		Model: i.config.Model,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	var response strings.Builder
	err := i.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if resp.Message.Content != "" {
			response.WriteString(resp.Message.Content)
			if onToken != nil {
				onToken(resp.Message.Content)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to query model: %w", err)
	}

	if onComplete != nil {
		onComplete(response.String())
	}

	return nil
}

// buildPrompt constructs an intelligent prompt based on context and providers
func (i *IntelSystem) buildPrompt(userQuery string, promptType PromptType) string {
	// Update context manager with current information
	i.updateContextManager(promptType)
	
	// Use context manager to build optimized prompt
	return i.contextManager.BuildPrompt(userQuery, promptType)
}

// updateContextManager updates the context manager with current information
func (i *IntelSystem) updateContextManager(promptType PromptType) {
	// 1. Add system prompt
	systemPrompt := i.config.SystemPrompt + `

RESPONSE STYLE GUIDELINES:
- Be concise and direct like Claude CLI
- Maximum 150 words total
- Use structured bullet points with ▸ 
- Include specific commands when relevant
- No verbose explanations or academic tone
- Focus on immediate, actionable steps
- Use professional CLI formatting
- Avoid redundant information
- End with clear next steps

FORMATTING RULES:
- Use ## for main sections
- Use ▸ for bullet points
- Use ` + "`code`" + ` for commands
- Keep paragraphs short (1-2 sentences)
- Use numbered lists for steps
- Avoid excessive technical jargon`

	i.contextManager.AddContext("system", ContextTypeSystem, systemPrompt, true)
	
	// 2. Add domain knowledge from providers
	for _, provider := range i.providers {
		knowledge := provider.GetDomainKnowledge()
		if knowledge != "" {
			i.contextManager.AddContext(
				fmt.Sprintf("domain-%s", provider.Name()),
				ContextTypeDomain,
				knowledge,
				true,
			)
		}
	}
	
	// 3. Add current state (only key information)
	for _, provider := range i.providers {
		state := provider.GetCurrentState()
		if len(state) > 0 {
			var stateStr strings.Builder
			stateStr.WriteString(fmt.Sprintf("Current %s:\n", provider.Name()))
			
			// Only include key state information relevant to current prompt type
			relevantKeys := i.getRelevantStateKeys(promptType)
			for _, key := range relevantKeys {
				if value, exists := state[key]; exists {
					stateStr.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
				}
			}
			
			i.contextManager.AddContext(
				fmt.Sprintf("state-%s", provider.Name()),
				ContextTypeState,
				stateStr.String(),
				false,
			)
		}
	}
	
	// 4. Add recent command history (last 3 only)
	i.context.mu.RLock()
	if len(i.context.RecentActions) > 0 {
		var historyStr strings.Builder
		historyStr.WriteString("Recent commands:\n")
		
		recentActions := i.context.RecentActions
		if len(recentActions) > 3 {
			recentActions = recentActions[len(recentActions)-3:]
		}
		
		for _, action := range recentActions {
			status := "✓"
			if !action.Success {
				status = "✗"
			}
			historyStr.WriteString(fmt.Sprintf("- %s %s %v\n", 
				status, action.Command, action.Args))
		}
		
		i.contextManager.AddContext(
			"history",
			ContextTypeHistory,
			historyStr.String(),
			false,
		)
	}
	i.context.mu.RUnlock()
	
	// 5. Add custom prompt template if available
	if template, exists := i.config.CustomPrompts[string(promptType)]; exists {
		promptStr := fmt.Sprintf("Task: %s", template)
		i.contextManager.AddContext(
			fmt.Sprintf("prompt-%s", promptType),
			ContextTypePrompt,
			promptStr,
			true,
		)
	}
	
	// Periodic cleanup
	i.contextManager.PruneHistory()
}

// GetOllamaStatus returns the current Ollama status
func (i *IntelSystem) GetOllamaStatus() (string, error) {
	return i.ollamaManager.GetStatus()
}

// GetOllamaManager returns the Ollama manager for advanced operations
func (i *IntelSystem) GetOllamaManager() *OllamaManager {
	return i.ollamaManager
}

// GetContextStats returns current context statistics
func (i *IntelSystem) GetContextStats() map[string]interface{} {
	return i.contextManager.GetStats()
}

// GetContextSummary returns a summary of current context
func (i *IntelSystem) GetContextSummary() string {
	return i.contextManager.GetContextSummary()
}

// ClearContext clears all context
func (i *IntelSystem) ClearContext() {
	i.contextManager.Clear()
}

// SetMaxTokens updates the maximum token limit
func (i *IntelSystem) SetMaxTokens(maxTokens int) {
	i.contextManager.SetMaxTokens(maxTokens)
}

// getRelevantStateKeys returns state keys relevant to the current prompt type
func (i *IntelSystem) getRelevantStateKeys(promptType PromptType) []string {
	switch promptType {
	case PromptAnalyze:
		return []string{"target_url", "authenticated", "total_discoveries", "high_severity_findings"}
	case PromptSuggest:
		return []string{"target_url", "authenticated", "discoveries_count", "schema_discovered"}
	case PromptExplain:
		return []string{"target_url", "authenticated"} // Minimal context for explanations
	default:
		return []string{"target_url", "authenticated", "total_discoveries"}
	}
}

// AnalyzeWithStreaming performs AI analysis with streaming output
func (i *IntelSystem) AnalyzeWithStreaming(userPrompt string) error {
	if !i.IsInitialized() {
		return fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(userPrompt, PromptAnalyze)
	
	// Use the regular query method and format the result
	content, err := i.queryModel(prompt)
	if err != nil {
		return err
	}
	
	// Format and display the response properly
	formatter := NewStreamingFormatter()
	formatter.FormatAndDisplayResponse(content)
	
	return nil
}

// SuggestWithStreaming provides AI-generated suggestions with streaming output
func (i *IntelSystem) SuggestWithStreaming(context string) error {
	if !i.IsInitialized() {
		return fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(context, PromptSuggest)
	
	// Use the regular query method and format the result
	content, err := i.queryModel(prompt)
	if err != nil {
		return err
	}
	
	// Format and display the response properly
	formatter := NewStreamingFormatter()
	formatter.FormatAndDisplayResponse(content)
	
	return nil
}

// ExplainWithStreaming provides detailed explanations with streaming output
func (i *IntelSystem) ExplainWithStreaming(topic string) error {
	if !i.IsInitialized() {
		return fmt.Errorf("Intel system not initialized")
	}

	prompt := i.buildPrompt(topic, PromptExplain)
	
	// Use the regular query method and format the result
	content, err := i.queryModel(prompt)
	if err != nil {
		return err
	}
	
	// Format and display the response properly
	formatter := NewStreamingFormatter()
	formatter.FormatAndDisplayResponse(content)
	
	return nil
}