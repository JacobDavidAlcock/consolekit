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
	Completions map[int]ArgumentCompletion // Argument completion configuration
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
		Completions: make(map[int]ArgumentCompletion),
	}
}

// RegisterFunc registers a function as a command handler
func (r *Registry) RegisterFunc(name string, fn func([]string) error, description string) {
	r.Register(name, HandlerFunc(fn), description)
}

// RegisterWithCompletion adds a command with custom completion support
func (r *Registry) RegisterWithCompletion(name string, handler Handler, description string, completions map[int]ArgumentCompletion) {
	r.commands[name] = &Command{
		Name:        name,
		Handler:     handler,
		Description: description,
		Subcommands: make(map[string]*Command),
		Completions: completions,
	}
}

// RegisterWithCompleter adds a command that implements the Completer interface
func (r *Registry) RegisterWithCompleter(name string, handler Handler, description string) {
	r.commands[name] = &Command{
		Name:        name,
		Handler:     handler,
		Description: description,
		Subcommands: make(map[string]*Command),
		Completions: make(map[int]ArgumentCompletion),
	}
}

// RegisterWithBuilder adds a command with completion built using CompletionBuilder
func (r *Registry) RegisterWithBuilder(name string, handler Handler, description string, builder *CompletionBuilder) {
	r.commands[name] = &Command{
		Name:        name,
		Handler:     handler,
		Description: description,
		Subcommands: make(map[string]*Command),
		Completions: builder.Build(),
	}
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

	for name, cmd := range r.commands {
		// Check if command implements Completer interface
		if completer, ok := cmd.Handler.(Completer); ok {
			items = append(items, r.buildDynamicCompletion(name, completer, cmd))
		} else if len(cmd.Completions) > 0 {
			// Use static completion configuration
			items = append(items, r.buildStaticCompletion(name, cmd))
		} else {
			// Fall back to hardcoded completions for known commands
			items = append(items, r.buildLegacyCompletion(name))
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

// buildDynamicCompletion creates completion for commands implementing Completer interface
func (r *Registry) buildDynamicCompletion(name string, completer Completer, cmd *Command) readline.PrefixCompleterInterface {
	// For dynamic completion, we start with a basic structure
	// The actual dynamic completion happens in the readline callback
	var subItems []readline.PrefixCompleterInterface
	
	// Get initial completions for empty args
	initialOptions := completer.Complete([]string{}, 0)
	for _, option := range initialOptions {
		subItems = append(subItems, readline.PcItem(option))
	}
	
	return readline.PcItem(name, subItems...)
}

// buildStaticCompletion creates completion from static configuration
func (r *Registry) buildStaticCompletion(name string, cmd *Command) readline.PrefixCompleterInterface {
	var subItems []readline.PrefixCompleterInterface
	
	// Add completions for each argument position
	for pos := 0; pos < 3; pos++ { // Support up to 3 argument positions
		if completion, exists := cmd.Completions[pos]; exists {
			// Add static options
			for _, option := range completion.Options {
				subItems = append(subItems, readline.PcItem(option))
			}
			
			// Add dynamic options if available
			if completion.Dynamic != nil {
				dynamicOptions := completion.Dynamic()
				for _, option := range dynamicOptions {
					subItems = append(subItems, readline.PcItem(option))
				}
			}
			
			// Add flag completions
			for flag, flagOptions := range completion.Flags {
				if len(flagOptions) > 0 {
					var flagSubItems []readline.PrefixCompleterInterface
					for _, flagOption := range flagOptions {
						flagSubItems = append(flagSubItems, readline.PcItem(flagOption))
					}
					subItems = append(subItems, readline.PcItem(flag, flagSubItems...))
				} else {
					subItems = append(subItems, readline.PcItem(flag))
				}
			}
		}
	}
	
	return readline.PcItem(name, subItems...)
}

// buildLegacyCompletion maintains backward compatibility for hardcoded completions
func (r *Registry) buildLegacyCompletion(name string) readline.PrefixCompleterInterface {
	switch name {
	case "intel":
		// Add Intel subcommands with proper nesting
		return readline.PcItem("intel",
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
		)
	case "show":
		// Add show subcommands with proper nesting
		return readline.PcItem("show",
			readline.PcItem("config"),
			readline.PcItem("session"),
			readline.PcItem("findings"),
		)
	case "set":
		// Add set subcommands with proper nesting
		return readline.PcItem("set",
			readline.PcItem("target"),
			readline.PcItem("timeout"),
			readline.PcItem("verbose"),
		)
	default:
		// Regular commands without subcommands
		return readline.PcItem(name)
	}
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