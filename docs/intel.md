# Intel AI Assistant

The Intel system provides AI-powered assistance for ConsoleKit applications using local language models via Ollama. It offers context-aware analysis, suggestions, and explanations tailored to your specific tool domain.

## Overview

Intel is designed to be:
- **Domain-agnostic**: Works with any CLI tool through pluggable context providers
- **Local**: Runs on your machine via Ollama, no cloud dependencies
- **Context-aware**: Understands your session state, command history, and discoveries
- **Extensible**: Easy to customize with domain-specific knowledge

## Quick Start

### Prerequisites

1. Install Ollama:
```bash
# macOS
brew install ollama

# Linux  
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama service
ollama serve
```

2. Add Intel to your application:
```go
import "github.com/jacobdavidalcock/consolekit/pkg/intel"

func main() {
    app := console.New("mytool")
    
    // One-line setup with domain expertise
    intel.QuickSetup(app, "mytool", "graphql", 
        "You are a GraphQL security expert...")
    
    app.Run()
}
```

3. Use Intel commands:
```bash
mytool > intel start
mytool > intel analyze
mytool > intel suggest next steps
mytool > intel explain GraphQL injection
```

## Architecture

### Core Components

- **IntelSystem**: Main coordinator handling LLM communication
- **ContextProvider**: Pluggable interface for domain-specific knowledge  
- **ModelManager**: Handles model selection and downloading
- **Commands**: Standard intel commands that work with any tool

### Context Providers

Context providers give Intel domain-specific knowledge about your tool:

```go
type ContextProvider interface {
    Name() string
    GetContext() (*ContextData, error)
    GetDomainKnowledge() string
    GetCurrentState() map[string]interface{}
    GetPromptTemplates() map[string]string
}
```

## Configuration

Configure Intel through YAML config files:

```yaml
intel:
  model: "phi3:3.8b"
  auto_download: true
  proactive: false
  context_depth: 10
  ollama_url: "http://localhost:11434"
  timeout: 30s
  
  custom_prompts:
    analyze: "Focus on security vulnerabilities and testing gaps"
    suggest: "Recommend specific next steps based on current findings"
    explain: "Provide detailed explanations with practical examples"
```

### Configuration Options

- `model`: LLM model to use (auto-selected by default)
- `auto_download`: Automatically download missing models
- `proactive`: Enable proactive suggestions (future feature)
- `context_depth`: Number of recent commands to include in context
- `ollama_url`: Ollama server URL
- `custom_prompts`: Override default prompts for different command types

## Model Selection

Intel automatically selects the best model based on:
- Available system RAM
- Tool domain (security, coding, general)
- User preferences

### Recommended Models

| Model | Size | Specialty | RAM Required | Description |
|-------|------|-----------|--------------|-------------|
| `phi3:3.8b` | 2.2GB | General | 4GB | Fast, capable general-purpose |
| `llama3.2:3b` | 2.0GB | General | 4GB | Excellent reasoning |
| `qwen2.5:3b` | 1.9GB | Coding | 4GB | Strong technical analysis |
| `gemma2:2b` | 1.6GB | Fast | 2GB | Lightweight option |

## Creating Context Providers

### Basic Provider

```go
type MyContextProvider struct {
    *intel.BaseContextProvider
    session *MySession
}

func NewMyContextProvider(session *MySession) *MyContextProvider {
    knowledge := `You are an expert in my domain. Key concepts:
    - Domain-specific concept 1
    - Domain-specific concept 2
    - Common vulnerabilities and patterns`
    
    base := intel.NewBaseContextProvider("my-context", "mydomain", knowledge)
    
    return &MyContextProvider{
        BaseContextProvider: base,
        session:            session,
    }
}

func (m *MyContextProvider) GetCurrentState() map[string]interface{} {
    return map[string]interface{}{
        "target":       m.session.Target,
        "authenticated": m.session.Authenticated,
        "findings":     len(m.session.Findings),
    }
}
```

### Advanced Provider with Custom Templates

```go
func (m *MyContextProvider) setupTemplates() {
    m.SetPromptTemplate(intel.PromptAnalyze, 
        "Analyze the current session focusing on domain-specific vulnerabilities...")
    
    m.SetPromptTemplate(intel.PromptSuggest, 
        "Based on current findings, suggest specific next steps...")
    
    m.SetPromptTemplate(intel.PromptExplain, 
        "Explain domain concepts with practical examples...")
}
```

## Integration Patterns

### Simple Integration

```go
// One-line setup for common use cases
intel.QuickSetup(app, "mytool", "security", domainKnowledge)
```

### Advanced Integration

```go
// Full control over configuration
config := &intel.Config{
    Model: "qwen2.5:3b",
    AutoDownload: true,
    CustomPrompts: map[string]string{
        "analyze": "Custom analysis prompt...",
    },
}

integration := intel.NewIntegration("mytool", config)
integration.WithProvider(myContextProvider)
integration.RegisterWith(app)
```

### Middleware Integration

```go
// Track command actions automatically
intelSystem := intel.New("mytool", config)

// In your command handlers
func (c *MyCommand) Execute(args []string) error {
    result := c.doWork(args)
    
    // Track for Intel context
    intelSystem.AddAction("mycommand", args, result, err == nil)
    
    return err
}
```

## Standard Commands

Intel provides these standard commands for any tool:

### `intel start`
Initializes the Intel system and downloads models if needed.

### `intel analyze [query]`
Analyzes current session state and provides insights. Optional query parameter for specific analysis.

### `intel suggest [context]`
Provides AI-generated suggestions for next steps based on current context.

### `intel explain <topic>`
Explains concepts, vulnerabilities, or techniques related to your domain.

### `intel status`
Shows Intel system status, active providers, and session information.

## Use Cases

### Security Testing Tools
- Analyze discovered vulnerabilities
- Suggest next testing steps  
- Explain attack vectors and payloads
- Recommend specific tools and techniques

### DevOps Tools
- Analyze configuration issues
- Suggest optimization strategies
- Explain infrastructure patterns
- Troubleshoot deployment problems

### Development Tools
- Code analysis and review
- Architecture suggestions
- Best practice recommendations
- Technology explanations

## Best Practices

### 1. Domain-Specific Knowledge
Provide comprehensive domain knowledge in your context provider:

```go
knowledge := `You are an expert in GraphQL security testing. Key concepts:
- Schema introspection and discovery
- Authorization bypass techniques
- Injection vulnerabilities in resolvers
- Rate limiting and DoS prevention
- Common GraphQL security tools and payloads`
```

### 2. Rich Context
Include relevant session state:

```go
func (p *MyProvider) GetCurrentState() map[string]interface{} {
    return map[string]interface{}{
        "target_url":     p.session.Target,
        "endpoints":      p.session.Endpoints,
        "vulnerabilities": p.session.Findings,
        "auth_status":    p.session.Authenticated,
        "scan_progress":  p.session.Progress,
    }
}
```

### 3. Custom Prompts
Tailor prompts for your specific use case:

```go
provider.SetPromptTemplate(intel.PromptSuggest, 
    "Based on the GraphQL schema and current findings, suggest specific "+
    "security tests including query examples and injection payloads.")
```

### 4. Action Tracking
Track important commands for context:

```go
// Track successful discoveries
if len(findings) > 0 {
    intel.AddAction("scan", args, fmt.Sprintf("Found %d issues", len(findings)), true)
}
```

## Troubleshooting

### Common Issues

**Intel not starting**
- Ensure Ollama is running: `ollama serve`
- Check Ollama URL in configuration
- Verify model availability: `ollama list`

**Model download failing**
- Check internet connection
- Ensure sufficient disk space (models are 1-4GB)
- Try manual download: `ollama pull phi3:3.8b`

**Poor AI responses**
- Improve domain knowledge in context provider
- Add more session context
- Use custom prompt templates
- Try a different model

**Performance issues**
- Use smaller models (gemma2:2b, llama3.2:1b)
- Reduce context depth in configuration
- Increase timeout for complex queries

### Getting Help

1. Check Intel system status: `intel status`
2. Review Ollama logs: `ollama logs`
3. Test with minimal context to isolate issues
4. Try different models to compare results

## Future Enhancements

- **Proactive suggestions**: Automatic recommendations during tool usage
- **Learning from sessions**: Improved suggestions based on user patterns  
- **Multi-model support**: Use different models for different types of queries
- **Cloud model support**: Integration with cloud-based LLM APIs
- **Voice interface**: Voice commands and responses
- **Plugin system**: Third-party Intel extensions