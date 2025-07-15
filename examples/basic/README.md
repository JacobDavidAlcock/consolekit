# Basic ConsoleKit Example

This example demonstrates the core features of ConsoleKit by creating a simple interactive CLI application.

## Running the Example

### Step 1: Install Dependencies
```bash
cd examples/basic
go mod tidy
```

### Step 2: Run the Example

#### Option 1: Run directly with go run
```bash
go run main.go
```

#### Option 2: Build and run the executable
```bash
go build -o example-cli
./example-cli
```

#### Option 3: Run with config file
```bash
go run main.go --config config.yaml
```

## Available Commands

### Original Commands
- `help` - Show all available commands
- `set <key> <value>` - Set a configuration value
- `show options` - Display all current configuration values
- `greet [name]` - Greet someone (defaults to "World")
- `demo` - Run a demonstration of progress indicators
- `make-config` - Generate an example configuration file
- `exit` or `quit` - Exit the application

### **NEW: Tab Completion Examples**
- `test` - **Dynamic completion** - Context-aware completion based on arguments
- `scan` - **Static completion** - Predefined completion using CompletionBuilder
- `pentest` - **Security pattern** - Security testing with predefined completion
- `request` - **HTTP pattern** - HTTP client with method and flag completion
- `analyze` - **Fluent interface** - Target analysis with fluent completion setup
- `connect` - **Dynamic external** - Service connection with external data completion

## Example Usage

```
$ go run main.go

┌─────────────────────────────────────────┐
│              Example CLI                │
│         Powered by ConsoleKit           │
└─────────────────────────────────────────┘

example > help

--- Example CLI Help Menu ---
  set                  Set a configuration variable
  show                 Display current configuration
  greet                Greet someone
  demo                 Run a demonstration
  make-config          Generate example configuration file
  exit / quit          Close the application.
  help                 Display this help menu.
------------------------

example > set username alice
[*] username => alice

example > set token secret123
[*] token => secret123

example > show options

--- Current State ---
  username        : alice
  token           : secr...123
--------------------

example > greet Alice
Hello, Alice! 👋

example > demo

🚀 Running demonstration...
[|] Processing data...
[/] Analyzing results...
✅ Demo completed! Found 5 items.

example > exit
```

## Features Demonstrated

### Core ConsoleKit Features
1. **Interactive Console**: Readline with command history and tab completion
2. **Command Registration**: Multiple types of command handlers
3. **State Management**: Thread-safe storage with sensitive data masking
4. **Progress Indicators**: Spinners and counters for long-running operations
5. **Configuration**: Support for config files and runtime settings
6. **Rich Output**: Colors, banners, and formatted display

### **NEW: Advanced Tab Completion**
7. **Dynamic Completion**: Commands implementing the `Completer` interface
8. **Static Completion**: Using `CompletionBuilder` for predefined options
9. **Quick Patterns**: Predefined security, HTTP, and file operation patterns
10. **Fluent Interface**: Readable command registration with completion
11. **External Data**: Dynamic completion from simulated external sources

## Testing Tab Completion

### Dynamic Completion Example
```bash
example > test <TAB>
# Shows: sql, xss, nosql, ldap, command, xxe, ssrf

example > test sql <TAB>  
# Shows: union, boolean, time, error (SQL-specific options)

example > test sql union --<TAB>
# Shows: --verbose, --quiet, --output, --format
```

### Static Completion Example
```bash
example > scan <TAB>
# Shows: endpoints, subdomains, directories, files

example > scan endpoints --<TAB>
# Shows: --threads, --timeout, --output, --verbose
```

### Security Pattern Example
```bash
example > pentest <TAB>
# Shows: security test types and predefined security flags

example > request <TAB>
# Shows: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
```

## Creating Your Own Commands

### Basic Command (No Completion)
To add a new command, implement the `command.Handler` interface:

```go
type MyCommand struct{}

func (c *MyCommand) Execute(args []string) error {
    fmt.Println("My custom command executed!")
    return nil
}

func (c *MyCommand) Description() string {
    return "My custom command description"
}

// Then register it:
app.AddCommand("mycmd", &MyCommand{}, "My custom command")
```

### **NEW: Command with Dynamic Completion**
Implement both `Handler` and `Completer` interfaces:

```go
type MyCommand struct{}

func (c *MyCommand) Execute(args []string) error {
    fmt.Printf("Executing with: %v\n", args)
    return nil
}

func (c *MyCommand) Description() string {
    return "Command with dynamic completion"
}

// Add dynamic completion
func (c *MyCommand) Complete(args []string, cursorPos int) []string {
    switch len(args) {
    case 0:
        return []string{"option1", "option2", "option3"}
    case 1:
        return []string{"--flag1", "--flag2", "--flag3"}
    default:
        return []string{}
    }
}

// Register with completer support:
app.AddCommandWithCompleter("mycmd", &MyCommand{}, "My command with completion")
```

### **NEW: Command with Static Completion**
Use `CompletionBuilder` for predefined options:

```go
builder := command.NewCompletionBuilder().
    AddPosition(0, "create", "read", "update", "delete").
    AddFlag("--format", "json", "yaml", "xml").
    AddFlag("--output", "stdout", "file").
    AddFlag("--verbose", "true", "false")

app.AddCommandWithBuilder("crud", &CrudCommand{}, "CRUD operations", builder)
```

### **NEW: Quick Pattern Registration**
Use predefined patterns for common use cases:

```go
ext := app.Extensions()
ext.SecurityTest("audit", &AuditCommand{}, "Security audit")
ext.HTTPClient("fetch", &FetchCommand{}, "HTTP client")
ext.FileOps("process", &ProcessCommand{}, "File processing")
```

### **NEW: Fluent Interface**
Use the fluent interface for readable registration:

```go
app.Fluent().
    Add("deploy", &DeployCommand{}).
    Desc("Deploy application").
    Arg(0, "staging", "production", "development").
    Flag("--rollback", "true", "false").
    Flag("--timeout", "30", "60", "120").
    Register()
```