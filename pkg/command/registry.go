package command

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

// Handler defines the interface for command handlers
type Handler interface {
	Execute(args []string) error
	Description() string
}

// HandlerFunc allows using functions as command handlers
type HandlerFunc func(args []string) error

func (f HandlerFunc) Execute(args []string) error {
	return f(args)
}

func (f HandlerFunc) Description() string {
	return "Custom command"
}

// Command represents a registered command
type Command struct {
	Name        string
	Handler     Handler
	Description string
	Subcommands map[string]*Command
}

// Registry manages command registration and execution
type Registry struct {
	commands map[string]*Command
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

// Register adds a new command to the registry
func (r *Registry) Register(name string, handler Handler, description string) {
	r.commands[name] = &Command{
		Name:        name,
		Handler:     handler,
		Description: description,
		Subcommands: make(map[string]*Command),
	}
}

// RegisterFunc registers a function as a command handler
func (r *Registry) RegisterFunc(name string, fn func([]string) error, description string) {
	r.Register(name, HandlerFunc(fn), description)
}

// Execute runs the specified command with arguments
func (r *Registry) Execute(name string, args []string) error {
	command, exists := r.commands[strings.ToLower(name)]
	if !exists {
		return fmt.Errorf("unknown command: %s. Type 'help' for a list of commands", name)
	}

	return command.Handler.Execute(args)
}

// BuildCompleter creates a readline completer from registered commands
func (r *Registry) BuildCompleter() readline.PrefixCompleterInterface {
	var items []readline.PrefixCompleterInterface

	for name := range r.commands {
		// Create subcommand completions for specific commands
		switch name {
		case "intel":
			// Add Intel subcommands with proper nesting
			items = append(items, readline.PcItem("intel",
				readline.PcItem("start"),
				readline.PcItem("analyze"),
				readline.PcItem("suggest"),
				readline.PcItem("explain"),
				readline.PcItem("status"),
				readline.PcItem("context",
					readline.PcItem("clear"),
					readline.PcItem("stats"),
					readline.PcItem("limit"),
				),
				readline.PcItem("validate",
					readline.PcItem("model"),
					readline.PcItem("url"),
					readline.PcItem("rules"),
				),
				readline.PcItem("help",
					readline.PcItem("errors"),
				),
			))
		case "show":
			// Add show subcommands with proper nesting
			items = append(items, readline.PcItem("show",
				readline.PcItem("config"),
				readline.PcItem("session"),
				readline.PcItem("findings"),
			))
		case "set":
			// Add set subcommands with proper nesting
			items = append(items, readline.PcItem("set",
				readline.PcItem("target"),
				readline.PcItem("timeout"),
				readline.PcItem("verbose"),
			))
		default:
			// Regular commands without subcommands
			items = append(items, readline.PcItem(name))
		}
	}

	// Add built-in commands
	items = append(items,
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)

	return readline.NewPrefixCompleter(items...)
}

// ShowHelp displays help for all registered commands
func (r *Registry) ShowHelp() {
	for name, cmd := range r.commands {
		fmt.Printf("  %-20s %s\n", name, cmd.Description)
	}
}

// GetCommand returns a command by name
func (r *Registry) GetCommand(name string) (*Command, bool) {
	cmd, exists := r.commands[strings.ToLower(name)]
	return cmd, exists
}

// ListCommands returns all registered command names
func (r *Registry) ListCommands() []string {
	var names []string
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}