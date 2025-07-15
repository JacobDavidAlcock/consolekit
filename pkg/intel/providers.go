package intel

import "time"

// ContextProvider defines the interface that tools must implement to provide domain-specific context
type ContextProvider interface {
	Name() string
	GetContext() (*ContextData, error)
	GetDomainKnowledge() string
	GetCurrentState() map[string]interface{}
	GetPromptTemplates() map[string]string
}

// ContextData represents the current context for AI analysis
type ContextData struct {
	Domain      string                 `json:"domain"`      // "firebase", "graphql", "kubernetes", etc.
	Session     map[string]interface{} `json:"session"`     // Current session data
	History     []Action               `json:"history"`     // Command history
	Discoveries []Finding              `json:"discoveries"` // What's been found/discovered
	State       map[string]interface{} `json:"state"`       // Current tool state
	Timestamp   time.Time              `json:"timestamp"`   // When context was created
}

// Action represents a command that was executed
type Action struct {
	Command     string                 `json:"command"`
	Args        []string               `json:"args"`
	Result      string                 `json:"result"`
	Success     bool                   `json:"success"`
	Timestamp   time.Time              `json:"timestamp"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// Finding represents something discovered during tool usage
type Finding struct {
	Type        string                 `json:"type"`        // "vulnerability", "endpoint", "data", etc.
	Severity    string                 `json:"severity"`    // "low", "medium", "high", "critical"
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Location    string                 `json:"location"`    // Where it was found
	Evidence    map[string]interface{} `json:"evidence"`    // Supporting data
	Timestamp   time.Time              `json:"timestamp"`
}

// PromptType defines different types of prompts for different use cases
type PromptType string

const (
	PromptAnalyze PromptType = "analyze"
	PromptSuggest PromptType = "suggest"
	PromptExplain PromptType = "explain"
	PromptDebug   PromptType = "debug"
	PromptHelp    PromptType = "help"
)

// BaseContextProvider provides a default implementation that tools can embed
type BaseContextProvider struct {
	name      string
	domain    string
	knowledge string
	templates map[string]string
}

// NewBaseContextProvider creates a new base context provider
func NewBaseContextProvider(name, domain, knowledge string) *BaseContextProvider {
	return &BaseContextProvider{
		name:      name,
		domain:    domain,
		knowledge: knowledge,
		templates: make(map[string]string),
	}
}

// Name returns the provider name
func (b *BaseContextProvider) Name() string {
	return b.name
}

// GetDomainKnowledge returns the domain-specific knowledge
func (b *BaseContextProvider) GetDomainKnowledge() string {
	return b.knowledge
}

// GetPromptTemplates returns the prompt templates
func (b *BaseContextProvider) GetPromptTemplates() map[string]string {
	return b.templates
}

// SetPromptTemplate sets a prompt template for a specific type
func (b *BaseContextProvider) SetPromptTemplate(promptType PromptType, template string) {
	b.templates[string(promptType)] = template
}

// GetContext provides a basic context implementation
func (b *BaseContextProvider) GetContext() (*ContextData, error) {
	return &ContextData{
		Domain:      b.domain,
		Session:     make(map[string]interface{}),
		History:     []Action{},
		Discoveries: []Finding{},
		State:       make(map[string]interface{}),
		Timestamp:   time.Now(),
	}, nil
}

// GetCurrentState provides a basic state implementation
func (b *BaseContextProvider) GetCurrentState() map[string]interface{} {
	return map[string]interface{}{
		"provider": b.name,
		"domain":   b.domain,
		"active":   true,
	}
}