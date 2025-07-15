package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jacobdavidalcock/consolekit/pkg/command"
	"github.com/jacobdavidalcock/consolekit/pkg/config"
	"github.com/jacobdavidalcock/consolekit/pkg/console"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// Example CLI application demonstrating ConsoleKit usage
func main() {
	// Create a new console application
	app := console.New("example-cli")
	app.WithPrompt("example > ")

	// Display a welcome banner
	banner := output.GenerateConsoleBanner("Example CLI", "Powered by ConsoleKit")
	app.SetBanner(banner)

	// Create application state
	state := config.NewState()

	// Register commands
	registerCommands(app, state)

	// Handle config file if provided
	if configPath, err := config.HandleStartupFlag(); err != nil {
		log.Fatal(err)
	} else if configPath != "" {
		cfg := config.New()
		if err := cfg.LoadFromFile(configPath); err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
		} else {
			fmt.Printf("✓ Config loaded from %s\n", configPath)
		}
	}

	// Start the interactive console
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// registerCommands sets up all available commands
func registerCommands(app *console.Console, state *config.State) {
	
	// SET command - mimics firescan's set functionality
	app.AddCommand("set", &SetCommand{state: state}, "Set a configuration variable")
	
	// SHOW command - mimics firescan's show functionality  
	app.AddCommand("show", &ShowCommand{state: state}, "Display current configuration")
	
	// GREET command - simple example command
	app.AddCommand("greet", command.HandlerFunc(func(args []string) error {
		name := "World"
		if len(args) > 0 {
			name = strings.Join(args, " ")
		}
		fmt.Printf("Hello, %s! 👋\n", name)
		return nil
	}), "Greet someone")
	
	// DEMO command - demonstrates progress indicators
	app.AddCommand("demo", &DemoCommand{}, "Run a demonstration")
	
	// CONFIG command - generate example config
	app.AddCommand("make-config", command.HandlerFunc(func(args []string) error {
		example := config.GenerateExample("example-cli")
		fmt.Println(example)
		return nil
	}), "Generate example configuration file")
	
	// Enhanced tab completion examples
	registerCompletionExamples(app, state)
}

// registerCompletionExamples demonstrates the new tab completion features
func registerCompletionExamples(app *console.Console, state *config.State) {
	// Example 1: Dynamic completion with Completer interface
	app.AddCommandWithCompleter("test", &TestCommand{state: state}, "Test with dynamic completion")
	
	// Example 2: Static completion with CompletionBuilder
	builder := command.NewCompletionBuilder().
		AddPosition(0, "endpoints", "subdomains", "directories", "files").
		AddFlag("--threads", "1", "5", "10", "20", "50").
		AddFlag("--timeout", "30", "60", "120", "300").
		AddFlag("--output", "json", "yaml", "text", "table").
		AddFlag("--verbose", "true", "false")
	
	app.AddCommandWithBuilder("scan", &ScanCommand{state: state}, "Scan with static completion", builder)
	
	// Example 3: Quick patterns for common use cases
	ext := app.Extensions()
	ext.SecurityTest("pentest", &PentestCommand{state: state}, "Security testing with predefined patterns")
	ext.HTTPClient("request", &RequestCommand{state: state}, "HTTP client with predefined patterns")
	
	// Example 4: Fluent interface
	app.Fluent().
		Add("analyze", &AnalyzeCommand{state: state}).
		Desc("Analyze target with fluent completion").
		Arg(0, "web", "api", "mobile", "network").
		Flag("--depth", "shallow", "medium", "deep").
		Flag("--format", "json", "xml", "yaml").
		SecurityFlags().
		Register()
	
	// Example 5: Dynamic external data
	app.Fluent().
		Add("connect", &ConnectCommand{state: state}).
		Desc("Connect to services with dynamic completion").
		DynamicArg(0, func() []string {
			// Simulate fetching from external source
			return []string{"database", "cache", "queue", "api"}
		}).
		DynamicFlag("--env", func() []string {
			// Simulate environment detection
			return []string{"development", "staging", "production"}
		}).
		Register()
}

// SetCommand handles setting configuration values
type SetCommand struct {
	state *config.State
}

func (c *SetCommand) Execute(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: set <key> <value>")
	}
	
	key := args[0]
	value := args[1]
	
	c.state.Set(key, value)
	fmt.Printf("[*] %s => %s\n", key, value)
	return nil
}

func (c *SetCommand) Description() string {
	return "Set a configuration variable"
}

// ShowCommand handles displaying configuration
type ShowCommand struct {
	state *config.State
}

func (c *ShowCommand) Execute(args []string) error {
	if len(args) == 0 || args[0] != "options" {
		return fmt.Errorf("usage: show options")
	}
	
	c.state.ShowAll()
	return nil
}

func (c *ShowCommand) Description() string {
	return "Display current configuration"
}

// DemoCommand demonstrates progress indicators
type DemoCommand struct{}

func (c *DemoCommand) Execute(args []string) error {
	fmt.Println("\n" + output.Yellow("🚀 Running demonstration..."))
	
	// Demonstrate spinner
	spinner := output.NewSpinner("Processing data...")
	spinner.Start()
	
	// Simulate work
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)
		if i == 15 {
			spinner.UpdateMessage("Analyzing results...")
		}
	}
	
	spinner.Stop()
	
	// Demonstrate progress counter
	counter := output.NewCounter("Scanning items", 50)
	counter.Start()
	
	// Simulate scanning work
	for i := 0; i < 50; i++ {
		time.Sleep(50 * time.Millisecond)
		counter.Increment()
		
		// Simulate finding something occasionally
		if i%10 == 0 && i > 0 {
			counter.IncrementFound()
		}
	}
	
	counter.Stop()
	
	fmt.Printf("\n%s✅ Demo completed! Found %d items.%s\n", 
		output.Green, counter.GetFound(), output.Reset)
	
	return nil
}

func (c *DemoCommand) Description() string {
	return "Run a demonstration of ConsoleKit features"
}

// New completion example commands

// TestCommand implements the Completer interface for dynamic completion
type TestCommand struct {
	state *config.State
}

func (c *TestCommand) Execute(args []string) error {
	fmt.Printf("🧪 Testing: %v\n", args)
	if len(args) > 0 {
		c.state.Set("last_test_type", args[0])
	}
	return nil
}

func (c *TestCommand) Description() string {
	return "Test with dynamic tab completion"
}

// Complete provides dynamic completion based on current context
func (c *TestCommand) Complete(args []string, cursorPos int) []string {
	switch len(args) {
	case 0:
		// First argument: test types
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
	default:
		// Additional flags
		return []string{"--verbose", "--quiet", "--output", "--format"}
	}
}

// ScanCommand uses static completion
type ScanCommand struct {
	state *config.State
}

func (c *ScanCommand) Execute(args []string) error {
	fmt.Printf("🔍 Scanning: %v\n", args)
	if len(args) > 0 {
		c.state.Set("last_scan_target", args[0])
	}
	return nil
}

func (c *ScanCommand) Description() string {
	return "Scan with static tab completion"
}

// PentestCommand uses security testing pattern
type PentestCommand struct {
	state *config.State
}

func (c *PentestCommand) Execute(args []string) error {
	fmt.Printf("🛡️ Penetration testing: %v\n", args)
	return nil
}

func (c *PentestCommand) Description() string {
	return "Penetration testing with predefined patterns"
}

// RequestCommand uses HTTP client pattern
type RequestCommand struct {
	state *config.State
}

func (c *RequestCommand) Execute(args []string) error {
	fmt.Printf("🌐 HTTP request: %v\n", args)
	return nil
}

func (c *RequestCommand) Description() string {
	return "HTTP client with predefined patterns"
}

// AnalyzeCommand uses fluent interface
type AnalyzeCommand struct {
	state *config.State
}

func (c *AnalyzeCommand) Execute(args []string) error {
	fmt.Printf("📊 Analyzing: %v\n", args)
	if len(args) > 0 {
		c.state.Set("last_analysis_type", args[0])
	}
	return nil
}

func (c *AnalyzeCommand) Description() string {
	return "Analyze with fluent interface completion"
}

// ConnectCommand uses dynamic external data
type ConnectCommand struct {
	state *config.State
}

func (c *ConnectCommand) Execute(args []string) error {
	fmt.Printf("🔌 Connecting: %v\n", args)
	if len(args) > 0 {
		c.state.Set("last_connection", args[0])
	}
	return nil
}

func (c *ConnectCommand) Description() string {
	return "Connect with dynamic completion from external data"
}