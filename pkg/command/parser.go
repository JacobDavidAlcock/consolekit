package command

import (
	"flag"
	"fmt"
	"strings"
)

// ParsedCommand represents a parsed command with its arguments and flags
type ParsedCommand struct {
	Name string
	Args []string
	Flags *flag.FlagSet
}

// Parser handles command line parsing and flag management
type Parser struct {
	flagSets map[string]*flag.FlagSet
}

// NewParser creates a new command parser
func NewParser() *Parser {
	return &Parser{
		flagSets: make(map[string]*flag.FlagSet),
	}
}

// ParseCommand parses a command line into command name and arguments
func ParseCommand(line string) (string, []string) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

// CreateFlagSet creates a new flag set for a command
func (p *Parser) CreateFlagSet(commandName string) *flag.FlagSet {
	flagSet := flag.NewFlagSet(commandName, flag.ContinueOnError)
	p.flagSets[commandName] = flagSet
	return flagSet
}

// GetFlagSet returns the flag set for a command
func (p *Parser) GetFlagSet(commandName string) (*flag.FlagSet, bool) {
	flagSet, exists := p.flagSets[commandName]
	return flagSet, exists
}

// ParseFlags parses command line flags for a specific command
func (p *Parser) ParseFlags(commandName string, args []string) error {
	flagSet, exists := p.flagSets[commandName]
	if !exists {
		return fmt.Errorf("no flag set defined for command: %s", commandName)
	}

	return flagSet.Parse(args)
}

// ValidateRequired checks that required flags are provided
func ValidateRequired(flagSet *flag.FlagSet, required []string) error {
	for _, name := range required {
		flag := flagSet.Lookup(name)
		if flag == nil {
			return fmt.Errorf("required flag not found: %s", name)
		}
		if flag.Value.String() == flag.DefValue {
			return fmt.Errorf("required flag not provided: --%s", name)
		}
	}
	return nil
}

// ShowUsage displays usage information for a command
func ShowUsage(commandName string, flagSet *flag.FlagSet) {
	fmt.Printf("Usage: %s [options]\n", commandName)
	fmt.Println("Options:")
	flagSet.PrintDefaults()
}