# Getting Started with ConsoleKit

ConsoleKit is a Go framework for building interactive command-line applications with minimal boilerplate. This guide will walk you through creating your first CLI application.

## Installation

```bash
go get github.com/jacobdavidalcock/consolekit
```

## Quick Start

### 1. Basic Console Application

Create a simple interactive console:

```go
package main

import (
    "github.com/jacobdavidalcock/consolekit/pkg/console"
)

func main() {
    app := console.New("myapp")
    app.Run()
}
```

This creates a basic console with:
- Interactive readline prompt (`myapp > `)
- Command history with arrow keys
- Tab completion for commands
- Built-in `help`, `exit`, and `quit` commands

### 2. Adding Commands

Register custom commands using the command registry:

```go
package main

import (
    "fmt"
    "github.com/jacobdavidalcock/consolekit/pkg/console"
    "github.com/jacobdavidalcock/consolekit/pkg/command"
)

func main() {
    app := console.New("myapp")
    
    // Add a simple command using HandlerFunc
    app.AddCommand("greet", command.HandlerFunc(func(args []string) error {
        name := "World"
        if len(args) > 0 {
            name = args[0]
        }
        fmt.Printf("Hello, %s!\n", name)
        return nil
    }), "Greet someone")
    
    app.Run()
}
```

### 3. Using State Management

ConsoleKit provides thread-safe state management for storing application data:

```go
package main

import (
    "fmt"
    "github.com/jacobdavidalcock/consolekit/pkg/console"
    "github.com/jacobdavidalcock/consolekit/pkg/config"
)

func main() {
    app := console.New("myapp")
    state := config.NewState()
    
    // Add a command that uses state
    app.AddCommand("set", &SetCommand{state: state}, "Set a value")
    app.AddCommand("get", &GetCommand{state: state}, "Get a value")
    
    app.Run()
}

type SetCommand struct {
    state *config.State
}

func (c *SetCommand) Execute(args []string) error {
    if len(args) != 2 {
        return fmt.Errorf("usage: set <key> <value>")
    }
    c.state.Set(args[0], args[1])
    fmt.Printf("Set %s = %s\n", args[0], args[1])
    return nil
}

func (c *SetCommand) Description() string {
    return "Set a configuration value"
}
```

### 4. Configuration Files

Load settings from YAML configuration files:

```go
package main

import (
    "fmt"
    "log"
    "github.com/jacobdavidalcock/consolekit/pkg/console"
    "github.com/jacobdavidalcock/consolekit/pkg/config"
)

func main() {
    app := console.New("myapp")
    
    // Handle --config flag
    if configPath, err := config.HandleStartupFlag(); err != nil {
        log.Fatal(err)
    } else if configPath != "" {
        cfg := config.New()
        if err := cfg.LoadFromFile(configPath); err != nil {
            fmt.Printf("Error loading config: %v\n", err)
        }
    }
    
    app.Run()
}
```

### 5. Rich Output and Progress

Add colored output, banners, and progress indicators:

```go
package main

import (
    "time"
    "github.com/jacobdavidalcock/consolekit/pkg/console"
    "github.com/jacobdavidalcock/consolekit/pkg/output"
    "github.com/jacobdavidalcock/consolekit/pkg/command"
)

func main() {
    app := console.New("myapp")
    
    // Set a banner
    banner := output.GenerateConsoleBanner("My App", "A demo application")
    app.SetBanner(banner)
    
    // Add a command with progress indicators
    app.AddCommand("work", command.HandlerFunc(func(args []string) error {
        spinner := output.NewSpinner("Processing...")
        spinner.Start()
        
        time.Sleep(2 * time.Second)
        
        spinner.Stop()
        fmt.Println(output.Green("✓ Work completed!"))
        return nil
    }), "Do some work")
    
    app.Run()
}
```

## Next Steps

- See [examples/basic/](../examples/basic/) for a complete working example
- Check the [API Reference](api.md) for detailed documentation
- Learn about [advanced features](advanced.md) like custom validation and complex command structures

## Common Patterns

### Command with Flags

```go
type ScanCommand struct{}

func (c *ScanCommand) Execute(args []string) error {
    parser := command.NewParser()
    flagSet := parser.CreateFlagSet("scan")
    
    verbose := flagSet.Bool("v", false, "Verbose output")
    target := flagSet.String("target", "", "Target to scan")
    
    if err := flagSet.Parse(args); err != nil {
        return err
    }
    
    if *target == "" {
        return fmt.Errorf("target is required")
    }
    
    if *verbose {
        fmt.Printf("Scanning target: %s\n", *target)
    }
    
    return nil
}
```

### Configuration Validation

```go
func setupValidation() *config.Validator {
    validator := config.NewValidator()
    
    // Email validation
    validator.AddEmailRule("email", true)
    
    // String with pattern
    validator.AddStringRule("projectid", true, `^[a-z0-9-]{6,30}$`)
    
    // Integer with range
    min, max := 1, 100
    validator.AddIntRule("port", false, &min, &max)
    
    return validator
}
```

This covers the basics of getting started with ConsoleKit. The framework handles all the interactive console complexity while letting you focus on your application logic.