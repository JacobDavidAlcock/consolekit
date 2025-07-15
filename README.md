# ConsoleKit

🚀 **A powerful Go framework for building interactive CLI applications with AI assistance**

ConsoleKit provides a complete toolkit for creating professional command-line interfaces with minimal boilerplate. Perfect for security tools, DevOps utilities, and any application that needs an interactive console with intelligent AI assistance.

## ✨ Key Features

### Core Framework
- **Interactive Console**: Full readline support with command history and tab completion
- **Command System**: Hierarchical command registration with automatic help generation
- **Configuration Management**: YAML config files with runtime updates and validation
- **State Management**: Thread-safe global state with secure credential display
- **Rich Output**: Colored output, progress indicators, JSON formatting, and ASCII banners

### Intel AI Assistant
- **Local LLM Integration**: Uses Ollama for privacy-focused AI assistance
- **Domain-Specific Context**: Pluggable context providers for specialized knowledge
- **Professional UI**: Claude-style responses with intelligent formatting
- **Smart Context Management**: Automatic context pruning and relevance scoring
- **Command Intelligence**: AI-powered analysis, suggestions, and explanations

## 🏗️ Architecture

```
pkg/
├── console/     # Core REPL and readline functionality
├── command/     # Command parsing, routing, and registration
├── config/      # Configuration loading, state management, validation
├── intel/       # AI assistant with local LLM integration
├── output/      # Colors, formatting, progress indicators
└── utils/       # Common utilities (strings, files, security)
```

## 🚀 Quick Start

### Installation
```bash
go get github.com/jacobdavidalcock/consolekit
```

### Simple Example
```go
package main

import "github.com/jacobdavidalcock/consolekit/pkg/console"

func main() {
    app := console.New("myapp")
    app.Run()
}
```

### Running Examples

#### Basic Example
```bash
cd examples/basic
go mod tidy
go run main.go
```

#### Intel AI Example (GraphQL Security Testing)
```bash
# Intel automatically installs and manages Ollama for you
cd examples/graphql-intel
go mod tidy
go run main.go

# Try Intel commands (Ollama will be installed automatically on first use):
graphql > intel start
graphql > intel analyze
graphql > intel explain graphql injection
graphql > intel suggest next steps
```

## 🤖 Intel AI Assistant

Intel provides Claude-quality AI assistance directly in your CLI applications:

### Quick Setup
```go
import "github.com/jacobdavidalcock/consolekit/pkg/intel"

func main() {
    app := console.New("mytool")
    
    // Enable AI assistant with domain expertise
    intel.QuickSetup(app, "mytool", "security", 
        "You are a security testing expert specializing in...")
    
    app.Run()
}
```

### Features
- **Zero Setup**: Automatically installs and manages Ollama - no manual setup required
- **Professional Responses**: Claude-style concise, actionable output
- **Smart Formatting**: ASCII art, proper markdown, and clean presentation
- **Context Awareness**: Understands your tool's current state and history
- **Tab Completion**: Full autocompletion for all intel commands
- **Intelligent Truncation**: Keeps responses CLI-appropriate with natural breaks
- **Domain Expertise**: Customizable knowledge for your specific use case

### Intel Commands
```bash
intel start              # Initialize AI system
intel analyze [query]    # Analyze current session
intel suggest [context]  # Get AI suggestions  
intel explain <topic>    # Detailed explanations
intel status            # System status
intel help              # Command reference
```

## 📖 Documentation

- [Getting Started Guide](docs/getting-started.md)
- [Intel AI Assistant](docs/intel.md)
- [API Reference](docs/api.md)
- [Examples](examples/)

## 🤝 Contributing

Contributions welcome! This framework was extracted from [FireScan](https://github.com/JacobDavidAlcock/firescan) to provide reusable CLI components.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.