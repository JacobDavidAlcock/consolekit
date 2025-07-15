# GraphQL Intel Example

This example demonstrates how to integrate the Intel AI assistant system with a GraphQL security testing tool using ConsoleKit.

## Prerequisites

1. **Go 1.21+** installed
2. **Ollama** running locally at `http://localhost:11434`
3. A compatible language model (will auto-download `phi3:3.8b` by default)

### Installing Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.ai/install.sh | sh

# Windows
# Download from https://ollama.ai/download

# Start Ollama service
ollama serve
```

## Running the Example

```bash
cd examples/graphql-intel
go mod tidy
go run main.go
```

## Intel Features Demonstrated

### 1. AI-Powered Analysis
```
graphql > intel start
🤖 Initializing Intel AI system...
✅ Intel AI system initialized successfully

graphql > intel analyze
🔍 Analyzing current session...
📊 Intel Analysis:
Based on your GraphQL testing session, I can see you haven't set a target yet. 
I recommend starting with introspection to discover the schema structure...
```

### 2. Smart Suggestions
```
graphql > target https://api.example.com/graphql
graphql > introspect
graphql > intel suggest

💡 Generating suggestions...
🎯 Intel Suggestions:
Next Steps:
1. Test authorization on the discovered user queries
2. Check for injection vulnerabilities in query parameters  
3. Test for recursive queries that could cause DoS
4. Verify proper error handling and information disclosure
```

### 3. Context-Aware Explanations
```
graphql > intel explain authorization bypass

📖 Explaining: authorization bypass...
📚 Intel Explanation:
GraphQL authorization bypass occurs when the API fails to properly validate 
permissions at the field or resolver level. Common attack vectors include:

- Direct object reference without permission checks
- Field-level authorization missing in resolvers  
- Token validation bypassed in certain query paths
```

## GraphQL-Specific Intelligence

The Intel system understands GraphQL security testing context:

- **Schema Analysis**: Analyzes discovered types, queries, and mutations
- **Vulnerability Context**: Understands GraphQL-specific attack vectors
- **Testing Guidance**: Provides specific testing steps for GraphQL APIs
- **Payload Suggestions**: Recommends GraphQL-specific payloads and techniques

## Example Session Flow

```bash
# 1. Start the application
go run main.go

# 2. Initialize AI assistant
graphql > intel start

# 3. Set target and discover schema
graphql > target https://api.example.com/graphql
graphql > introspect

# 4. Get AI analysis of discoveries
graphql > intel analyze

# 5. Run security scans
graphql > scan

# 6. Get AI suggestions for next steps
graphql > intel suggest

# 7. Ask for explanations of findings
graphql > intel explain authorization bypass

# 8. Check session status
graphql > show
graphql > intel status
```

## Custom Context Provider

The example shows how to create a domain-specific context provider:

```go
// GraphQLContextProvider provides GraphQL-specific context
type GraphQLContextProvider struct {
    *intel.BaseContextProvider
    session *GraphQLSession
}

// Provides domain knowledge about GraphQL security
func (g *GraphQLContextProvider) GetDomainKnowledge() string {
    return `GraphQL security concepts: introspection, authorization, 
    injection vulnerabilities, recursive queries...`
}

// Provides current session state for AI analysis
func (g *GraphQLContextProvider) GetCurrentState() map[string]interface{} {
    return map[string]interface{}{
        "target_url": g.session.Target,
        "authenticated": g.session.Authenticated,
        "schema_discovered": len(g.session.Schema) > 0,
        "total_discoveries": len(g.session.Discoveries),
    }
}
```

## Configuration

The Intel system can be configured through YAML:

```yaml
intel:
  model: "phi3:3.8b"
  auto_download: true
  proactive: false
  context_depth: 10
  ollama_url: "http://localhost:11434"
  
  custom_prompts:
    analyze: "Focus on GraphQL security vulnerabilities and testing gaps"
    suggest: "Recommend specific GraphQL security testing techniques"
    explain: "Provide detailed GraphQL security explanations with examples"
```

## Available Commands

### Core GraphQL Commands
- `target <url>` - Set GraphQL endpoint
- `introspect` - Discover schema via introspection
- `query <graphql>` - Execute GraphQL queries
- `scan` - Run automated security scans
- `auth <token>` - Set authentication token
- `show` - Display session information

### Intel AI Commands
- `intel start` - Initialize AI assistant
- `intel analyze` - Get AI analysis of current session
- `intel suggest` - Get AI suggestions for next steps  
- `intel explain <topic>` - Get detailed explanations
- `intel status` - Show Intel system status

## Model Recommendations

For GraphQL security testing:
- **phi3:3.8b** (default) - Good general knowledge, fast
- **qwen2.5:3b** - Strong coding/technical analysis
- **llama3.2:3b** - Excellent reasoning capabilities

The system will automatically select the best model based on your system resources.

## Integration with Other Tools

This pattern can be easily adapted for other security tools:

```go
// Enable Intel for any ConsoleKit app
intel.QuickSetup(app, "mytool", "domain", domainKnowledge)

// Or create custom integration
integration := intel.NewIntegration("mytool", config)
integration.WithProvider(myContextProvider)
integration.RegisterWith(app)
```