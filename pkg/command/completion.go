package command

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// Completer interface allows commands to provide custom tab completion
type Completer interface {
	Complete(args []string, cursorPos int) []string
}

// ArgumentCompletion defines completion options for command arguments
type ArgumentCompletion struct {
	Position int                    // Argument position (0-based)
	Options  []string              // Static completion options
	Dynamic  func() []string       // Dynamic option generator
	Flags    map[string][]string   // Flag completions (--flag -> options)
}

// CompletionContext provides context for completion
type CompletionContext struct {
	Command   string   // Current command name
	Args      []string // Current arguments
	CursorPos int      // Position in the argument list
}

// CompletionResult contains completion suggestions
type CompletionResult struct {
	Suggestions []string // List of completion suggestions
	HasMore     bool     // Whether there are more options available
}

// StaticCompletion creates completion from static options
func StaticCompletion(options ...string) func() []string {
	return func() []string {
		return options
	}
}

// FileCompletion provides file path completion
func FileCompletion(extensions ...string) func() []string {
	return func() []string {
		// Simple file completion - in a real implementation this would
		// integrate with readline's file completion
		return []string{"file1.txt", "file2.json", "config.yaml"}
	}
}

// NumberCompletion provides number range completion
func NumberCompletion(min, max, step int) func() []string {
	return func() []string {
		var options []string
		for i := min; i <= max; i += step {
			options = append(options, strconv.Itoa(i))
		}
		return options
	}
}

// BoolCompletion provides boolean completion
func BoolCompletion() func() []string {
	return func() []string {
		return []string{"true", "false", "yes", "no", "on", "off"}
	}
}

// CompletionBuilder helps build completion configurations
type CompletionBuilder struct {
	completions map[int]ArgumentCompletion
	flags       map[string][]string
}

// NewCompletionBuilder creates a new completion builder
func NewCompletionBuilder() *CompletionBuilder {
	return &CompletionBuilder{
		completions: make(map[int]ArgumentCompletion),
		flags:       make(map[string][]string),
	}
}

// AddPosition adds completion for a specific argument position
func (cb *CompletionBuilder) AddPosition(pos int, options ...string) *CompletionBuilder {
	cb.completions[pos] = ArgumentCompletion{
		Position: pos,
		Options:  options,
	}
	return cb
}

// AddDynamicPosition adds dynamic completion for a specific argument position
func (cb *CompletionBuilder) AddDynamicPosition(pos int, generator func() []string) *CompletionBuilder {
	cb.completions[pos] = ArgumentCompletion{
		Position: pos,
		Dynamic:  generator,
	}
	return cb
}

// AddFlag adds completion for a flag
func (cb *CompletionBuilder) AddFlag(flag string, options ...string) *CompletionBuilder {
	cb.flags[flag] = options
	return cb
}

// AddDynamicFlag adds dynamic completion for a flag
func (cb *CompletionBuilder) AddDynamicFlag(flag string, generator func() []string) *CompletionBuilder {
	cb.flags[flag] = generator()
	return cb
}

// Build creates the final completion configuration
func (cb *CompletionBuilder) Build() map[int]ArgumentCompletion {
	// Apply flags to all completions
	for pos, completion := range cb.completions {
		if completion.Flags == nil {
			completion.Flags = make(map[string][]string)
		}
		for flag, options := range cb.flags {
			completion.Flags[flag] = options
		}
		cb.completions[pos] = completion
	}
	return cb.completions
}

// DefaultCompletion provides common completion patterns
type DefaultCompletion struct{}

// GetCommonFlags returns commonly used flags
func (dc *DefaultCompletion) GetCommonFlags() map[string][]string {
	return map[string][]string{
		"--help":     {},
		"--verbose":  BoolCompletion()(),
		"--quiet":    BoolCompletion()(),
		"--output":   {"json", "yaml", "text", "table"},
		"--format":   {"json", "yaml", "xml", "csv", "table"},
		"--timeout":  NumberCompletion(1, 300, 5)(),
		"--threads":  NumberCompletion(1, 100, 1)(),
		"--config":   {"config.yaml", "settings.json"},
	}
}

// GetSecurityTestTypes returns common security test types
func (dc *DefaultCompletion) GetSecurityTestTypes() []string {
	return []string{
		"sql", "nosql", "xss", "xxe", "ssrf", "lfi", "rfi",
		"command", "ldap", "xpath", "template", "ssti",
	}
}

// GetHTTPMethods returns HTTP method completions
func (dc *DefaultCompletion) GetHTTPMethods() []string {
	return []string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS",
	}
}

// GetOutputFormats returns common output formats
func (dc *DefaultCompletion) GetOutputFormats() []string {
	return []string{
		"json", "yaml", "xml", "csv", "table", "text", "html",
	}
}

// CompletionHelper provides utilities for building completions
type CompletionHelper struct {
	defaults *DefaultCompletion
}

// NewCompletionHelper creates a new completion helper
func NewCompletionHelper() *CompletionHelper {
	return &CompletionHelper{
		defaults: &DefaultCompletion{},
	}
}

// ParseCurrentArgument extracts the current argument being completed
func (ch *CompletionHelper) ParseCurrentArgument(args []string, cursorPos int) (string, bool) {
	if cursorPos >= len(args) {
		return "", false
	}
	return args[cursorPos], true
}

// IsFlag checks if the current argument is a flag
func (ch *CompletionHelper) IsFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

// ExtractFlagName extracts the flag name from an argument
func (ch *CompletionHelper) ExtractFlagName(arg string) string {
	if strings.HasPrefix(arg, "--") {
		return strings.Split(arg, "=")[0]
	}
	if strings.HasPrefix(arg, "-") {
		return arg[:2] // Short flag like -v
	}
	return ""
}

// FilterByPrefix filters options by a given prefix
func (ch *CompletionHelper) FilterByPrefix(options []string, prefix string) []string {
	var filtered []string
	for _, option := range options {
		if strings.HasPrefix(strings.ToLower(option), strings.ToLower(prefix)) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

// BuildPrefixCompleter creates a readline PrefixCompleter from argument completions
func (ch *CompletionHelper) BuildPrefixCompleter(commandName string, completions map[int]ArgumentCompletion) readline.PrefixCompleterInterface {
	// This is a simplified implementation
	// In practice, this would need to work with readline's dynamic completion
	var items []readline.PrefixCompleterInterface
	
	// Add static completions for first argument
	if completion, exists := completions[0]; exists {
		for _, option := range completion.Options {
			items = append(items, readline.PcItem(option))
		}
		
		// Add flag completions
		for flag := range completion.Flags {
			items = append(items, readline.PcItem(flag))
		}
	}
	
	return readline.PcItem(commandName, items...)
}

// GetFileExtensions extracts file extensions from a path pattern
func (ch *CompletionHelper) GetFileExtensions(pattern string) []string {
	if strings.Contains(pattern, "*") {
		ext := filepath.Ext(pattern)
		if ext != "" {
			return []string{ext}
		}
	}
	return []string{}
}

// MergeCompletions merges multiple completion sources
func (ch *CompletionHelper) MergeCompletions(completions ...[]string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, completion := range completions {
		for _, item := range completion {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
	}
	
	return result
}