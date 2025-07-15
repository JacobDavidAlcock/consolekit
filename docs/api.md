# ConsoleKit API Reference

This document provides detailed API documentation for all ConsoleKit packages.

## Package: console

The `console` package provides the main interactive console functionality.

### type Console

```go
type Console struct {
    Name         string
    Prompt       string
    HistoryFile  string
    Commands     *command.Registry
}
```

The main console application struct.

#### func New

```go
func New(name string) *Console
```

Creates a new Console instance with default settings.

#### func (*Console) WithPrompt

```go
func (c *Console) WithPrompt(prompt string) *Console
```

Sets a custom prompt string. Returns the console for method chaining.

#### func (*Console) WithHistoryFile

```go
func (c *Console) WithHistoryFile(file string) *Console
```

Sets a custom history file location. Returns the console for method chaining.

#### func (*Console) AddCommand

```go
func (c *Console) AddCommand(name string, handler command.Handler, description string)
```

Registers a new command with the console.

#### func (*Console) SetBanner

```go
func (c *Console) SetBanner(banner string)
```

Displays a startup banner when the console starts.

#### func (*Console) Run

```go
func (c *Console) Run() error
```

Starts the interactive console REPL loop. This is a blocking call that runs until the user exits.

#### func (*Console) Close

```go
func (c *Console) Close() error
```

Gracefully shuts down the console and cleans up resources.

## Package: command

The `command` package handles command registration, parsing, and execution.

### type Handler

```go
type Handler interface {
    Execute(args []string) error
    Description() string
}
```

Interface that all command handlers must implement.

### type HandlerFunc

```go
type HandlerFunc func(args []string) error
```

Adapter to allow using functions as command handlers.

#### func (HandlerFunc) Execute

```go
func (f HandlerFunc) Execute(args []string) error
```

#### func (HandlerFunc) Description

```go
func (f HandlerFunc) Description() string
```

### type Registry

```go
type Registry struct {
    // contains filtered or unexported fields
}
```

Manages command registration and execution.

#### func NewRegistry

```go
func NewRegistry() *Registry
```

Creates a new command registry.

#### func (*Registry) Register

```go
func (r *Registry) Register(name string, handler Handler, description string)
```

Adds a new command to the registry.

#### func (*Registry) RegisterFunc

```go
func (r *Registry) RegisterFunc(name string, fn func([]string) error, description string)
```

Registers a function as a command handler.

#### func (*Registry) Execute

```go
func (r *Registry) Execute(name string, args []string) error
```

Runs the specified command with arguments.

### type Parser

```go
type Parser struct {
    // contains filtered or unexported fields
}
```

Handles command line parsing and flag management.

#### func NewParser

```go
func NewParser() *Parser
```

Creates a new command parser.

#### func (*Parser) CreateFlagSet

```go
func (p *Parser) CreateFlagSet(commandName string) *flag.FlagSet
```

Creates a new flag set for a command.

#### func (*Parser) ParseFlags

```go
func (p *Parser) ParseFlags(commandName string, args []string) error
```

Parses command line flags for a specific command.

## Package: config

The `config` package provides configuration management and state handling.

### type Config

```go
type Config struct {
    // contains filtered or unexported fields
}
```

Represents application configuration.

#### func New

```go
func New() *Config
```

Creates a new config instance.

#### func (*Config) LoadFromFile

```go
func (c *Config) LoadFromFile(path string) error
```

Loads configuration from a YAML file.

#### func (*Config) SaveToFile

```go
func (c *Config) SaveToFile(path string) error
```

Saves configuration to a YAML file.

#### func (*Config) Set

```go
func (c *Config) Set(key string, value interface{})
```

Sets a configuration value.

#### func (*Config) Get

```go
func (c *Config) Get(key string) (interface{}, bool)
```

Gets a configuration value.

#### func (*Config) GetString

```go
func (c *Config) GetString(key string) (string, bool)
```

Gets a string configuration value.

### type State

```go
type State struct {
    // contains filtered or unexported fields
}
```

Manages global application state with thread safety.

#### func NewState

```go
func NewState() *State
```

Creates a new state manager.

#### func (*State) Set

```go
func (s *State) Set(key string, value interface{})
```

Sets a state value (thread-safe).

#### func (*State) Get

```go
func (s *State) Get(key string) (interface{}, bool)
```

Gets a state value (thread-safe).

#### func (*State) GetString

```go
func (s *State) GetString(key string) (string, bool)
```

Gets a string state value (thread-safe).

#### func (*State) ShowAll

```go
func (s *State) ShowAll()
```

Displays all state values with sensitive data masked.

### type Validator

```go
type Validator struct {
    // contains filtered or unexported fields
}
```

Defines validation rules for configuration values.

#### func NewValidator

```go
func NewValidator() *Validator
```

Creates a new validator.

#### func (*Validator) AddStringRule

```go
func (v *Validator) AddStringRule(key string, required bool, pattern string) error
```

Adds a string validation rule with optional regex pattern.

#### func (*Validator) AddEmailRule

```go
func (v *Validator) AddEmailRule(key string, required bool)
```

Adds an email validation rule.

#### func (*Validator) Validate

```go
func (v *Validator) Validate(config map[string]interface{}) error
```

Validates a configuration map against defined rules.

## Package: output

The `output` package provides formatting, colors, and progress indicators.

### Color Constants

```go
const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Cyan   = "\033[36m"
    Bold   = "\033[1m"
)
```

#### func Colorize

```go
func Colorize(text, color string) string
```

Wraps text with the specified color.

#### func Red

```go
func Red(text string) string
```

Returns red colored text.

#### func Green

```go
func Green(text string) string
```

Returns green colored text.

### type ProgressSpinner

```go
type ProgressSpinner struct {
    // contains filtered or unexported fields
}
```

Displays a spinning progress indicator.

#### func NewSpinner

```go
func NewSpinner(message string) *ProgressSpinner
```

Creates a new progress spinner.

#### func (*ProgressSpinner) Start

```go
func (p *ProgressSpinner) Start()
```

Begins the spinner animation.

#### func (*ProgressSpinner) Stop

```go
func (p *ProgressSpinner) Stop()
```

Stops the spinner animation.

### type ProgressCounter

```go
type ProgressCounter struct {
    // contains filtered or unexported fields
}
```

Displays progress with counters.

#### func NewCounter

```go
func NewCounter(message string, total int64) *ProgressCounter
```

Creates a new progress counter.

#### func (*ProgressCounter) Increment

```go
func (p *ProgressCounter) Increment()
```

Increments the current counter.

#### func PrintBanner

```go
func PrintBanner(banner, color string)
```

Displays an ASCII art banner with optional color.

#### func GenerateConsoleBanner

```go
func GenerateConsoleBanner(appName, description string) string
```

Creates a banner for console applications.

## Package: utils

The `utils` package provides common utility functions.

#### func MaskString

```go
func MaskString(s string, prefixLen, suffixLen int) string
```

Hides the middle of a string for secure display (extracted from firescan).

#### func GenerateCaseVariations

```go
func GenerateCaseVariations(word string) []string
```

Takes a word and returns lowercase, PascalCase, and UPPERCASE variations (extracted from firescan).

#### func LoadWordlist

```go
func LoadWordlist(filePath string) ([]string, error)
```

Loads a wordlist from a file, returning a slice of words (extracted from firescan).

#### func FileExists

```go
func FileExists(filePath string) bool
```

Checks if a file exists.

#### func GetConfigDir

```go
func GetConfigDir(appName string) (string, error)
```

Returns a configuration directory path.

## Package: intel

The `intel` package provides AI assistant functionality with local LLM integration.

### func QuickSetup

```go
func QuickSetup(app *console.Console, appName, domain, knowledge string)
```

One-line setup for Intel AI assistant with domain expertise.

### type IntelSystem

```go
type IntelSystem struct {
    // contains filtered or unexported fields
}
```

Core AI assistant system coordinator.

### type ContextProvider

```go
type ContextProvider interface {
    Name() string
    GetContext() (*ContextData, error)
    GetDomainKnowledge() string
    GetCurrentState() map[string]interface{}
    GetPromptTemplates() map[string]string
}
```

Interface for providing domain-specific context to the AI system.

### type Config

```go
type Config struct {
    Model         string
    AutoDownload  bool
    Proactive     bool
    ContextDepth  int
    SystemPrompt  string
    CustomPrompts map[string]string
    OllamaURL     string
    Timeout       time.Duration
}
```

Configuration for the Intel AI system.

### Standard Commands

Intel automatically registers these commands:
- `intel start` - Initialize the AI system
- `intel analyze [query]` - Analyze current session
- `intel suggest [context]` - Get AI suggestions
- `intel explain <topic>` - Detailed explanations
- `intel status` - System status
- `intel help` - Command reference

For detailed Intel documentation, see the [Intel AI Guide](intel.md).

---

This API reference covers the main public interfaces of ConsoleKit. For usage examples, see the [Getting Started Guide](getting-started.md).