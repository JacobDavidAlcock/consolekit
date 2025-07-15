package command

// Quick provides convenient helper functions for common completion patterns
type Quick struct{}

// SecurityTest creates completion for security testing commands
func (q *Quick) SecurityTest() *CompletionBuilder {
	defaults := &DefaultCompletion{}
	return NewCompletionBuilder().
		AddPosition(0, defaults.GetSecurityTestTypes()...).
		AddFlag("--target", "http://example.com", "https://api.example.com").
		AddFlag("--threads", "1", "5", "10", "20", "50").
		AddFlag("--timeout", "30", "60", "120", "300").
		AddFlag("--output", defaults.GetOutputFormats()...).
		AddFlag("--verbose", "true", "false").
		AddFlag("--quiet", "true", "false")
}

// Discovery creates completion for discovery/scanning commands
func (q *Quick) Discovery() *CompletionBuilder {
	defaults := &DefaultCompletion{}
	return NewCompletionBuilder().
		AddPosition(0, "endpoints", "subdomains", "directories", "files", "parameters").
		AddFlag("--target", "http://example.com", "https://api.example.com").
		AddFlag("--wordlist", "common.txt", "directories.txt", "files.txt").
		AddFlag("--threads", NumberCompletion(1, 100, 5)()...).
		AddFlag("--delay", "0", "100", "500", "1000").
		AddFlag("--output", defaults.GetOutputFormats()...).
		AddFlag("--extensions", "php", "asp", "jsp", "html", "js").
		AddFlag("--status-codes", "200", "301", "302", "403", "500")
}

// HTTPClient creates completion for HTTP client commands
func (q *Quick) HTTPClient() *CompletionBuilder {
	defaults := &DefaultCompletion{}
	return NewCompletionBuilder().
		AddPosition(0, defaults.GetHTTPMethods()...).
		AddPosition(1, "http://example.com", "https://api.example.com").
		AddFlag("--method", defaults.GetHTTPMethods()...).
		AddFlag("--header", "Content-Type: application/json", "Authorization: Bearer token").
		AddFlag("--timeout", "30", "60", "120").
		AddFlag("--follow-redirects", BoolCompletion()()...).
		AddFlag("--verify-ssl", BoolCompletion()()...).
		AddFlag("--output", defaults.GetOutputFormats()...)
}

// FileOperations creates completion for file-related commands
func (q *Quick) FileOperations() *CompletionBuilder {
	return NewCompletionBuilder().
		AddDynamicPosition(0, FileCompletion(".txt", ".json", ".yaml", ".xml")).
		AddFlag("--format", "json", "yaml", "xml", "csv").
		AddFlag("--output", "output.txt", "results.json", "report.html").
		AddFlag("--append", BoolCompletion()()...).
		AddFlag("--backup", BoolCompletion()()...)
}

// Configuration creates completion for configuration commands
func (q *Quick) Configuration() *CompletionBuilder {
	return NewCompletionBuilder().
		AddPosition(0, "get", "set", "list", "reset", "export", "import").
		AddPosition(1, "target", "timeout", "threads", "output", "format", "verbose").
		AddFlag("--config", "config.yaml", "settings.json").
		AddFlag("--global", BoolCompletion()()...).
		AddFlag("--user", BoolCompletion()()...)
}

// Reporting creates completion for reporting commands
func (q *Quick) Reporting() *CompletionBuilder {
	defaults := &DefaultCompletion{}
	return NewCompletionBuilder().
		AddPosition(0, "generate", "export", "view", "summary").
		AddFlag("--format", defaults.GetOutputFormats()...).
		AddFlag("--template", "default", "detailed", "executive", "technical").
		AddFlag("--include", "summary", "details", "recommendations", "appendix").
		AddFlag("--exclude", "debug", "verbose", "raw-data").
		AddFlag("--output", "report.html", "report.pdf", "report.json")
}

// Database creates completion for database-related commands
func (q *Quick) Database() *CompletionBuilder {
	return NewCompletionBuilder().
		AddPosition(0, "connect", "query", "export", "import", "schema").
		AddFlag("--host", "localhost", "127.0.0.1").
		AddFlag("--port", "3306", "5432", "1433", "27017").
		AddFlag("--database", "mysql", "postgresql", "mssql", "mongodb").
		AddFlag("--username", "root", "admin", "user").
		AddFlag("--ssl", BoolCompletion()()...).
		AddFlag("--timeout", "30", "60", "120")
}

// Registry extension methods for easier registration
type RegistryExtensions struct {
	registry *Registry
}

// NewRegistryExtensions creates extensions for a registry
func NewRegistryExtensions(registry *Registry) *RegistryExtensions {
	return &RegistryExtensions{registry: registry}
}

// SecurityTest registers a security testing command
func (re *RegistryExtensions) SecurityTest(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.SecurityTest())
}

// Discovery registers a discovery/scanning command
func (re *RegistryExtensions) Discovery(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.Discovery())
}

// HTTPClient registers an HTTP client command
func (re *RegistryExtensions) HTTPClient(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.HTTPClient())
}

// FileOps registers a file operations command
func (re *RegistryExtensions) FileOps(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.FileOperations())
}

// Config registers a configuration command
func (re *RegistryExtensions) Config(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.Configuration())
}

// Report registers a reporting command
func (re *RegistryExtensions) Report(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.Reporting())
}

// Database registers a database command
func (re *RegistryExtensions) Database(name string, handler Handler, description string) {
	quick := &Quick{}
	re.registry.RegisterWithBuilder(name, handler, description, quick.Database())
}

// Custom allows registering with a custom completion builder
func (re *RegistryExtensions) Custom(name string, handler Handler, description string, builder *CompletionBuilder) {
	re.registry.RegisterWithBuilder(name, handler, description, builder)
}

// Fluent interface for the Registry
type FluentRegistry struct {
	registry *Registry
	quick    *Quick
}

// NewFluentRegistry creates a fluent interface for command registration
func NewFluentRegistry(registry *Registry) *FluentRegistry {
	return &FluentRegistry{
		registry: registry,
		quick:    &Quick{},
	}
}

// Add starts a fluent command registration
func (fr *FluentRegistry) Add(name string, handler Handler) *FluentCommand {
	return &FluentCommand{
		registry: fr.registry,
		name:     name,
		handler:  handler,
		builder:  NewCompletionBuilder(),
	}
}

// FluentCommand allows fluent-style command configuration
type FluentCommand struct {
	registry    *Registry
	name        string
	handler     Handler
	description string
	builder     *CompletionBuilder
}

// Desc sets the command description
func (fc *FluentCommand) Desc(description string) *FluentCommand {
	fc.description = description
	return fc
}

// Arg adds argument completion at a specific position
func (fc *FluentCommand) Arg(position int, options ...string) *FluentCommand {
	fc.builder.AddPosition(position, options...)
	return fc
}

// DynamicArg adds dynamic argument completion
func (fc *FluentCommand) DynamicArg(position int, generator func() []string) *FluentCommand {
	fc.builder.AddDynamicPosition(position, generator)
	return fc
}

// Flag adds flag completion
func (fc *FluentCommand) Flag(flag string, options ...string) *FluentCommand {
	fc.builder.AddFlag(flag, options...)
	return fc
}

// DynamicFlag adds dynamic flag completion
func (fc *FluentCommand) DynamicFlag(flag string, generator func() []string) *FluentCommand {
	fc.builder.AddDynamicFlag(flag, generator)
	return fc
}

// CommonFlags adds commonly used flags
func (fc *FluentCommand) CommonFlags() *FluentCommand {
	defaults := &DefaultCompletion{}
	flags := defaults.GetCommonFlags()
	for flag, options := range flags {
		fc.builder.AddFlag(flag, options...)
	}
	return fc
}

// SecurityFlags adds security testing flags
func (fc *FluentCommand) SecurityFlags() *FluentCommand {
	fc.Flag("--target", "http://example.com", "https://api.example.com")
	fc.Flag("--threads", NumberCompletion(1, 50, 5)()...)
	fc.Flag("--timeout", "30", "60", "120", "300")
	fc.Flag("--output", "json", "yaml", "text", "table")
	return fc
}

// HTTPFlags adds HTTP-related flags
func (fc *FluentCommand) HTTPFlags() *FluentCommand {
	defaults := &DefaultCompletion{}
	fc.Flag("--method", defaults.GetHTTPMethods()...)
	fc.Flag("--timeout", "30", "60", "120")
	fc.Flag("--follow-redirects", BoolCompletion()()...)
	fc.Flag("--verify-ssl", BoolCompletion()()...)
	return fc
}

// Register completes the fluent registration
func (fc *FluentCommand) Register() {
	fc.registry.RegisterWithBuilder(fc.name, fc.handler, fc.description, fc.builder)
}

// Quick methods for common patterns
func (fc *FluentCommand) AsSecurityTest() {
	quick := &Quick{}
	fc.registry.RegisterWithBuilder(fc.name, fc.handler, fc.description, quick.SecurityTest())
}

func (fc *FluentCommand) AsDiscovery() {
	quick := &Quick{}
	fc.registry.RegisterWithBuilder(fc.name, fc.handler, fc.description, quick.Discovery())
}

func (fc *FluentCommand) AsHTTPClient() {
	quick := &Quick{}
	fc.registry.RegisterWithBuilder(fc.name, fc.handler, fc.description, quick.HTTPClient())
}