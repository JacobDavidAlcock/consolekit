package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jacobdavidalcock/consolekit/pkg/command"
	"github.com/jacobdavidalcock/consolekit/pkg/console"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

func main() {
	app := console.New("completion-demo")
	app.WithPrompt("demo > ")

	banner := output.GenerateConsoleBanner("Completion Demo", "Advanced Tab Completion with ConsoleKit")
	app.SetBanner(banner)

	// Register commands with different completion patterns
	registerCommands(app)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerCommands(app *console.Console) {
	// Example 1: Command implementing Completer interface
	app.AddCommandWithCompleter("test", &TestCommand{}, "Test security vulnerabilities with dynamic completion")

	// Example 2: Using CompletionBuilder for static completion
	builder := command.NewCompletionBuilder().
		AddPosition(0, "endpoints", "subdomains", "directories", "files").
		AddFlag("--threads", "1", "5", "10", "20", "50").
		AddFlag("--timeout", "30", "60", "120", "300").
		AddFlag("--output", "json", "yaml", "text", "table").
		AddFlag("--verbose", "true", "false")

	app.AddCommandWithBuilder("discover", &DiscoverCommand{}, "Discover resources", builder)

	// Example 3: Using registry extensions for quick patterns
	ext := app.Extensions()
	ext.SecurityTest("scan", &ScanCommand{}, "Security scan with predefined completion")
	ext.HTTPClient("request", &RequestCommand{}, "HTTP client with predefined completion")

	// Example 4: Using fluent interface
	app.Fluent().
		Add("analyze", &AnalyzeCommand{}).
		Desc("Analyze target with custom completion").
		Arg(0, "web", "api", "mobile", "network").
		Flag("--depth", "shallow", "medium", "deep").
		Flag("--format", "json", "xml", "yaml").
		SecurityFlags().
		Register()

	// Example 5: Dynamic completion with external data
	app.Fluent().
		Add("connect", &ConnectCommand{}).
		Desc("Connect to external services").
		DynamicArg(0, func() []string {
			// This could fetch from database, config file, etc.
			return []string{"database", "api", "cache", "queue"}
		}).
		DynamicFlag("--env", func() []string {
			// This could read from environment or config
			return []string{"development", "staging", "production"}
		}).
		Register()

	// Example 6: File operations with file completion
	quick := &command.Quick{}
	app.AddCommandWithBuilder("export", &ExportCommand{}, "Export data to file", quick.FileOperations())

	// Example 7: Database operations
	app.AddCommandWithBuilder("db", &DatabaseCommand{}, "Database operations", quick.Database())
}

// TestCommand implements the Completer interface for dynamic completion
type TestCommand struct{}

func (c *TestCommand) Execute(args []string) error {
	fmt.Printf("🔍 Running test: %v\n", args)
	return nil
}

func (c *TestCommand) Description() string {
	return "Test security vulnerabilities with dynamic completion"
}

// Complete provides dynamic completion based on current context
func (c *TestCommand) Complete(args []string, cursorPos int) []string {
	switch len(args) {
	case 0:
		// First argument: injection types
		return []string{"sql", "xss", "nosql", "ldap", "command", "xxe", "ssrf"}
	case 1:
		// Second argument: depends on first argument
		switch args[0] {
		case "sql":
			return []string{"union", "boolean", "time", "error"}
		case "xss":
			return []string{"reflected", "stored", "dom"}
		case "nosql":
			return []string{"mongodb", "couchdb", "redis"}
		default:
			return []string{"--target", "--threads", "--timeout"}
		}
	case 2:
		// Handle flags
		if strings.HasPrefix(args[1], "--") {
			switch args[1] {
			case "--threads":
				return []string{"1", "5", "10", "20", "50"}
			case "--timeout":
				return []string{"30", "60", "120", "300"}
			case "--target":
				return []string{"http://example.com", "https://api.example.com"}
			}
		}
		return []string{"--verbose", "--quiet", "--output"}
	default:
		// Additional flags
		return []string{"--verbose", "--quiet", "--output", "--format"}
	}
}

// Simple command implementations for demonstration
type DiscoverCommand struct{}

func (c *DiscoverCommand) Execute(args []string) error {
	fmt.Printf("🔍 Discovering: %v\n", args)
	return nil
}

func (c *DiscoverCommand) Description() string {
	return "Discover resources with static completion"
}

type ScanCommand struct{}

func (c *ScanCommand) Execute(args []string) error {
	fmt.Printf("🛡️ Scanning: %v\n", args)
	return nil
}

func (c *ScanCommand) Description() string {
	return "Security scan with predefined completion patterns"
}

type RequestCommand struct{}

func (c *RequestCommand) Execute(args []string) error {
	fmt.Printf("🌐 Making request: %v\n", args)
	return nil
}

func (c *RequestCommand) Description() string {
	return "HTTP client with predefined completion patterns"
}

type AnalyzeCommand struct{}

func (c *AnalyzeCommand) Execute(args []string) error {
	fmt.Printf("📊 Analyzing: %v\n", args)
	return nil
}

func (c *AnalyzeCommand) Description() string {
	return "Analyze target with fluent interface completion"
}

type ConnectCommand struct{}

func (c *ConnectCommand) Execute(args []string) error {
	fmt.Printf("🔌 Connecting: %v\n", args)
	return nil
}

func (c *ConnectCommand) Description() string {
	return "Connect to services with dynamic completion"
}

type ExportCommand struct{}

func (c *ExportCommand) Execute(args []string) error {
	fmt.Printf("📤 Exporting: %v\n", args)
	return nil
}

func (c *ExportCommand) Description() string {
	return "Export data with file completion"
}

type DatabaseCommand struct{}

func (c *DatabaseCommand) Execute(args []string) error {
	fmt.Printf("🗄️ Database operation: %v\n", args)
	return nil
}

func (c *DatabaseCommand) Description() string {
	return "Database operations with predefined completion"
}