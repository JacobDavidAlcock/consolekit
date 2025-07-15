package console

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jacobdavidalcock/consolekit/pkg/command"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// Console represents the main interactive console
type Console struct {
	Name         string
	Prompt       string
	HistoryFile  string
	Commands     *command.Registry
	readline     *readline.Instance
}

// New creates a new Console instance
func New(name string) *Console {
	return &Console{
		Name:        name,
		Prompt:      name + " > ",
		HistoryFile: "/tmp/" + name + "_history.tmp",
		Commands:    command.NewRegistry(),
	}
}

// WithPrompt sets a custom prompt
func (c *Console) WithPrompt(prompt string) *Console {
	c.Prompt = prompt
	return c
}

// WithHistoryFile sets a custom history file location
func (c *Console) WithHistoryFile(file string) *Console {
	c.HistoryFile = file
	return c
}

// AddCommand registers a new command
func (c *Console) AddCommand(name string, handler command.Handler, description string) {
	c.Commands.Register(name, handler, description)
}

// AddCompleter adds tab completion for a command
func (c *Console) AddCompleter(completer readline.PrefixCompleterInterface) {
	// This will be set when initializing readline
}

// EnableIntel adds Intel AI assistant capabilities to the console
func (c *Console) EnableIntel(intel interface{}) {
	// This will register Intel commands when the intel package is imported
	// The interface{} type avoids circular dependencies
	if intelSystem, ok := intel.(interface {
		RegisterCommands(*Console)
	}); ok {
		intelSystem.RegisterCommands(c)
	}
}

// SetBanner displays a startup banner
func (c *Console) SetBanner(banner string) {
	fmt.Println(output.Cyan(banner))
}

// Run starts the interactive console REPL
func (c *Console) Run() error {
	// Create completer from registered commands
	completer := c.Commands.BuildCompleter()

	// Setup readline configuration
	config := &readline.Config{
		Prompt:          c.Prompt,
		HistoryFile:     c.HistoryFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	c.readline = rl

	// Main REPL loop (extracted from firescan)
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt || err == io.EOF {
			break
		}

		input := strings.Fields(line)
		if len(input) == 0 {
			continue
		}

		commandName := input[0]
		args := input[1:]

		// Handle built-in commands
		switch strings.ToLower(commandName) {
		case "exit", "quit":
			return nil
		case "help":
			c.showHelp()
			continue
		}

		// Execute registered command
		if err := c.Commands.Execute(commandName, args); err != nil {
			fmt.Printf("❌ %s\n", err.Error())
		}
	}

	return nil
}

// showHelp displays help for all registered commands
func (c *Console) showHelp() {
	fmt.Printf("\n--- %s Help Menu ---\n", c.Name)
	c.Commands.ShowHelp()
	fmt.Println("  exit / quit           Close the application.")
	fmt.Println("  help                  Display this help menu.")
	fmt.Println("------------------------")
}

// Close gracefully shuts down the console
func (c *Console) Close() error {
	if c.readline != nil {
		return c.readline.Close()
	}
	return nil
}