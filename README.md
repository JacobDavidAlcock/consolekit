# ConsoleKit

🚀 **A powerful Go framework for building interactive CLI applications**

ConsoleKit provides a complete toolkit for creating professional command-line interfaces with minimal boilerplate. Perfect for security tools, DevOps utilities, and any application that needs an interactive console.

## ✨ Features

- **Interactive Console**: Full readline support with command history and tab completion
- **Command System**: Hierarchical command registration with automatic help generation
- **Configuration Management**: YAML config files with runtime updates and validation
- **State Management**: Thread-safe global state with secure credential display
- **Rich Output**: Colored output, progress indicators, JSON formatting, and ASCII banners
- **Intel AI Assistant**: Local LLM integration with domain-specific context providers
- **Utilities**: String masking, file operations, input validation, and more

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
# Requires Ollama running locally
cd examples/graphql-intel
go mod tidy
go run main.go

# Try Intel commands:
# intel start
# intel analyze  
# intel suggest
```

### Intel AI Assistant

Add AI capabilities to your CLI with one line:

```go
import "github.com/jacobdavidalcock/consolekit/pkg/intel"

func main() {
    app := console.New("mytool")
    
    // Enable AI assistant with domain expertise
    intel.QuickSetup(app, "mytool", "security", 
        "You are a security testing expert...")
    
    app.Run()
}
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