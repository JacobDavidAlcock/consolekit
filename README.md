# ConsoleKit

🚀 **A powerful Go framework for building interactive CLI applications**

ConsoleKit provides a complete toolkit for creating professional command-line interfaces with minimal boilerplate. Perfect for security tools, DevOps utilities, and any application that needs an interactive console.

## ✨ Features

- **Interactive Console**: Full readline support with command history and tab completion
- **Command System**: Hierarchical command registration with automatic help generation
- **Configuration Management**: YAML config files with runtime updates and validation
- **State Management**: Thread-safe global state with secure credential display
- **Rich Output**: Colored output, progress indicators, JSON formatting, and ASCII banners
- **Utilities**: String masking, file operations, input validation, and more

## 🏗️ Architecture

```
pkg/
├── console/     # Core REPL and readline functionality
├── command/     # Command parsing, routing, and registration
├── config/      # Configuration loading, state management, validation
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

### Running the Example
```bash
# Navigate to the basic example
cd examples/basic

# Install dependencies
go mod tidy

# Run directly
go run main.go

# Or build and run
go build -o example-cli
./example-cli
```

## 📖 Documentation

- [Getting Started Guide](docs/getting-started.md)
- [API Reference](docs/api.md)
- [Examples](examples/)

## 🤝 Contributing

Contributions welcome! This framework was extracted from [FireScan](https://github.com/JacobDavidAlcock/firescan) to provide reusable CLI components.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.