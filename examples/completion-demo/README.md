# Tab Completion Demo

This example demonstrates the enhanced tab completion system in ConsoleKit.

## Features Demonstrated

### 1. Dynamic Completion with Completer Interface
```go
type TestCommand struct{}

func (c *TestCommand) Complete(args []string, cursorPos int) []string {
    // Custom completion logic based on current arguments
}
```

**Try:** `test <TAB>` to see injection types, then `test sql <TAB>` for SQL-specific options.

### 2. Static Completion with CompletionBuilder
```go
builder := command.NewCompletionBuilder().
    AddPosition(0, "endpoints", "subdomains", "directories", "files").
    AddFlag("--threads", "1", "5", "10", "20", "50")
```

**Try:** `discover <TAB>` for discovery types, `discover endpoints --<TAB>` for flags.

### 3. Quick Registration Patterns
```go
ext := app.Extensions()
ext.SecurityTest("scan", &ScanCommand{}, "Security scan")
ext.HTTPClient("request", &RequestCommand{}, "HTTP client")
```

**Try:** `scan <TAB>` or `request <TAB>` for predefined security/HTTP completions.

### 4. Fluent Interface
```go
app.Fluent().
    Add("analyze", &AnalyzeCommand{}).
    Desc("Analyze target").
    Arg(0, "web", "api", "mobile", "network").
    Flag("--depth", "shallow", "medium", "deep").
    SecurityFlags().
    Register()
```

**Try:** `analyze <TAB>` for target types, `analyze web --<TAB>` for flags.

### 5. Dynamic External Data
```go
DynamicArg(0, func() []string {
    // Fetch from database, config, etc.
    return []string{"database", "api", "cache", "queue"}
})
```

**Try:** `connect <TAB>` for dynamically generated service types.

## Running the Demo

```bash
cd examples/completion-demo
go mod tidy
go run main.go
```

## Available Commands

- `test` - Dynamic completion based on injection types
- `discover` - Static completion for discovery operations  
- `scan` - Security testing with predefined patterns
- `request` - HTTP client with predefined patterns
- `analyze` - Fluent interface with custom completion
- `connect` - Dynamic completion from external data
- `export` - File operations with file completion
- `db` - Database operations with predefined patterns

## Tab Completion Tips

1. **Press TAB** after command names to see available arguments
2. **Press TAB** after `--` to see available flags
3. **Type partial text** then TAB to filter options
4. **Use multiple arguments** - completion changes based on context

## Architecture

The completion system supports:

- **Static completion** - Predefined options
- **Dynamic completion** - Generated at runtime
- **Context-aware completion** - Changes based on previous arguments
- **Flag completion** - Specific options for each flag
- **Nested completion** - Multi-level argument structures
- **External data** - Integration with databases, APIs, configs

This provides a flexible foundation for building intelligent CLI tools with excellent user experience.