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