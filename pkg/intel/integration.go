package intel

import (
	"github.com/jacobdavidalcock/consolekit/pkg/console"
)

// Integration provides easy Intel integration for ConsoleKit applications
type Integration struct {
	system *IntelSystem
}

// NewIntegration creates a new Intel integration
func NewIntegration(appName string, config *Config) *Integration {
	return &Integration{
		system: New(appName, config),
	}
}

// WithProvider adds a context provider to the Intel system
func (i *Integration) WithProvider(provider ContextProvider) *Integration {
	i.system.RegisterProvider(provider)
	return i
}

// WithConfig updates the Intel system configuration
func (i *Integration) WithConfig(config *Config) *Integration {
	i.system.config = config
	return i
}

// RegisterWith adds Intel commands to a ConsoleKit application
func (i *Integration) RegisterWith(app *console.Console) {
	RegisterIntelCommands(app, i.system)
}

// GetSystem returns the underlying Intel system for advanced usage
func (i *Integration) GetSystem() *IntelSystem {
	return i.system
}

// Helper function to create a quick Intel integration
func EnableIntel(app *console.Console, appName string, providers ...ContextProvider) *IntelSystem {
	config := DefaultConfig()
	intel := New(appName, config)
	
	for _, provider := range providers {
		intel.RegisterProvider(provider)
	}
	
	RegisterIntelCommands(app, intel)
	return intel
}

// QuickSetup provides a one-line Intel setup for common use cases
func QuickSetup(app *console.Console, appName, domain, knowledge string) *IntelSystem {
	// Create a basic context provider
	provider := NewBaseContextProvider(appName+"-context", domain, knowledge)
	
	// Set up domain-specific prompt templates
	switch domain {
	case "graphql":
		provider.SetPromptTemplate(PromptAnalyze, 
			"Analyze GraphQL session. List 3-5 key findings with severity. Focus on schema, auth, injection risks. Suggest immediate next steps.")
		provider.SetPromptTemplate(PromptSuggest, 
			"Suggest 3-5 GraphQL security commands to run next. Include specific queries and mutations to test. Use code examples.")
		provider.SetPromptTemplate(PromptExplain, 
			"Explain GraphQL concept briefly. Include 2-3 security risks and practical test commands. Keep under 100 words.")
	case "firebase":
		provider.SetPromptTemplate(PromptAnalyze, 
			"Analyze Firebase session. List database rules, auth issues, and data exposure risks. Suggest immediate fixes.")
		provider.SetPromptTemplate(PromptSuggest, 
			"Suggest 3-5 Firebase security tests to run next. Include specific database queries and auth bypasses.")
		provider.SetPromptTemplate(PromptExplain, 
			"Explain Firebase concept briefly. Include security risks and test commands. Keep under 100 words.")
	case "kubernetes", "k8s":
		provider.SetPromptTemplate(PromptAnalyze, 
			"Analyze K8s session. List RBAC, pod security, and cluster misconfigurations. Suggest immediate actions.")
		provider.SetPromptTemplate(PromptSuggest, 
			"Suggest 3-5 Kubernetes security commands to run next. Include kubectl commands and security checks.")
		provider.SetPromptTemplate(PromptExplain, 
			"Explain K8s concept briefly. Include security risks and kubectl commands. Keep under 100 words.")
	}
	
	return EnableIntel(app, appName, provider)
}

// ActionMiddleware can be used to automatically track command actions for Intel context
func (i *Integration) ActionMiddleware(command string, args []string, result string, success bool) {
	i.system.AddAction(command, args, result, success)
}