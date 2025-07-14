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

Once the application starts, you can use these commands:

- `help` - Show all available commands
- `set <key> <value>` - Set a configuration value
- `show options` - Display all current configuration values
- `greet [name]` - Greet someone (defaults to "World")
- `demo` - Run a demonstration of progress indicators
- `make-config` - Generate an example configuration file
- `exit` or `quit` - Exit the application

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

1. **Interactive Console**: Readline with command history and tab completion
2. **Command Registration**: Multiple types of command handlers
3. **State Management**: Thread-safe storage with sensitive data masking
4. **Progress Indicators**: Spinners and counters for long-running operations
5. **Configuration**: Support for config files and runtime settings
6. **Rich Output**: Colors, banners, and formatted display

## Creating Your Own Commands

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

Or use a simple function:

```go
app.AddCommand("simple", command.HandlerFunc(func(args []string) error {
    fmt.Println("Simple command!")
    return nil
}), "A simple command")
```