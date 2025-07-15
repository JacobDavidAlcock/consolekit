package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jacobdavidalcock/consolekit/pkg/config"
	"github.com/jacobdavidalcock/consolekit/pkg/console"
	"github.com/jacobdavidalcock/consolekit/pkg/intel"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// GraphQLSession represents the current GraphQL testing session
type GraphQLSession struct {
	Target       string
	Schema       map[string]interface{}
	Endpoints    []string
	Queries      []string
	Mutations    []string
	Discoveries  []intel.Finding
	Authenticated bool
	Token        string
}

// GraphQLContextProvider provides GraphQL-specific context to Intel
type GraphQLContextProvider struct {
	*intel.BaseContextProvider
	session *GraphQLSession
}

// NewGraphQLContextProvider creates a GraphQL context provider
func NewGraphQLContextProvider(session *GraphQLSession) *GraphQLContextProvider {
	knowledge := `GraphQL Security Testing:

COMMON VULNERABILITIES:
- Introspection enabled in production
- Missing authorization checks
- SQL/NoSQL injection in resolvers
- Recursive queries (DoS)
- Information disclosure via errors
- IDOR vulnerabilities

TEST COMMANDS:
- Introspection: {__schema{types{name}}}
- Auth bypass: Try mutations without auth
- Injection: Test variables with payloads
- DoS: Deep nested queries
- Error leakage: Invalid field names`

	base := intel.NewBaseContextProvider("graphql-context", "graphql", knowledge)
	
	// Set GraphQL-specific prompt templates
	base.SetPromptTemplate(intel.PromptAnalyze, 
		"Analyze the GraphQL security testing session. Review discovered endpoints, schema information, and potential vulnerabilities. Suggest specific attack vectors based on the current findings.")
	
	base.SetPromptTemplate(intel.PromptSuggest, 
		"Based on the GraphQL testing session, suggest specific next steps. Include recommended queries, mutations to test, authorization bypass techniques, and potential injection points.")
	
	base.SetPromptTemplate(intel.PromptExplain, 
		"Explain GraphQL security concepts, vulnerabilities, or testing techniques. Provide practical examples and specific payloads when relevant.")

	return &GraphQLContextProvider{
		BaseContextProvider: base,
		session:            session,
	}
}

// GetContext provides current GraphQL session context
func (g *GraphQLContextProvider) GetContext() (*intel.ContextData, error) {
	context := &intel.ContextData{
		Domain:      "graphql",
		Session:     make(map[string]interface{}),
		History:     []intel.Action{},
		Discoveries: g.session.Discoveries,
		State:       g.GetCurrentState(),
		Timestamp:   time.Now(),
	}

	// Add session-specific data
	context.Session["target"] = g.session.Target
	context.Session["authenticated"] = g.session.Authenticated
	context.Session["endpoints_count"] = len(g.session.Endpoints)
	context.Session["queries_found"] = len(g.session.Queries)
	context.Session["mutations_found"] = len(g.session.Mutations)
	context.Session["discoveries_count"] = len(g.session.Discoveries)

	return context, nil
}

// GetCurrentState returns current session state
func (g *GraphQLContextProvider) GetCurrentState() map[string]interface{} {
	state := g.BaseContextProvider.GetCurrentState()
	
	state["target_url"] = g.session.Target
	state["authenticated"] = g.session.Authenticated
	state["schema_discovered"] = len(g.session.Schema) > 0
	state["endpoints"] = g.session.Endpoints
	state["total_discoveries"] = len(g.session.Discoveries)
	
	if len(g.session.Discoveries) > 0 {
		var highSeverity, mediumSeverity, lowSeverity int
		for _, finding := range g.session.Discoveries {
			switch finding.Severity {
			case "high", "critical":
				highSeverity++
			case "medium":
				mediumSeverity++
			case "low", "info":
				lowSeverity++
			}
		}
		state["high_severity_findings"] = highSeverity
		state["medium_severity_findings"] = mediumSeverity
		state["low_severity_findings"] = lowSeverity
	}

	return state
}

func main() {
	// Create console application
	app := console.New("graphqlstrike")
	app.WithPrompt("graphql > ")

	// Initialize session
	session := &GraphQLSession{
		Endpoints:   []string{},
		Queries:     []string{},
		Mutations:   []string{},
		Discoveries: []intel.Finding{},
	}

	// Create state for configuration
	state := config.NewState()

	// Set banner
	banner := output.GenerateConsoleBanner("GraphQLStrike", "AI-Powered GraphQL Security Testing")
	app.SetBanner(banner)

	// Register core commands
	registerCommands(app, session, state)

	// Set up Intel AI assistant with GraphQL expertise
	setupIntelligence(app, session)

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

// setupIntelligence configures the Intel AI system
func setupIntelligence(app *console.Console, session *GraphQLSession) {
	// Create GraphQL context provider
	provider := NewGraphQLContextProvider(session)

	// Set up Intel with the context provider
	intel.QuickSetup(app, "graphqlstrike", "graphql", provider.GetDomainKnowledge())

	fmt.Printf("%sIntel AI assistant ready! Try: intel start%s\n", output.CyanColor, output.Reset)
}

// registerCommands sets up all application commands
func registerCommands(app *console.Console, session *GraphQLSession, state *config.State) {
	
	// TARGET command - set GraphQL endpoint
	app.AddCommand("target", &TargetCommand{session: session, state: state}, "Set GraphQL target endpoint")
	
	// INTROSPECT command - discover GraphQL schema
	app.AddCommand("introspect", &IntrospectCommand{session: session}, "Discover GraphQL schema via introspection")
	
	// QUERY command - execute GraphQL queries
	app.AddCommand("query", &QueryCommand{session: session}, "Execute GraphQL queries")
	
	// SCAN command - automated vulnerability scanning
	app.AddCommand("scan", &ScanCommand{session: session}, "Run automated GraphQL security scans")
	
	// SHOW command - display session information
	app.AddCommand("show", &ShowCommand{session: session, state: state}, "Display session information")
	
	// AUTH command - authentication management
	app.AddCommand("auth", &AuthCommand{session: session}, "Manage authentication")
}

// TargetCommand handles setting the GraphQL endpoint
type TargetCommand struct {
	session *GraphQLSession
	state   *config.State
}

func (c *TargetCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: target <url>")
	}
	
	url := args[0]
	c.session.Target = url
	c.state.Set("target", url)
	
	fmt.Printf("Target set to: %s%s%s\n", output.CyanColor, url, output.Reset)
	fmt.Printf("Try: %sintrospect%s to discover the schema\n", output.YellowColor, output.Reset)
	return nil
}

func (c *TargetCommand) Description() string {
	return "Set GraphQL target endpoint"
}

// IntrospectCommand handles GraphQL schema discovery
type IntrospectCommand struct {
	session *GraphQLSession
}

func (c *IntrospectCommand) Execute(args []string) error {
	if c.session.Target == "" {
		return fmt.Errorf("no target set. Use 'target <url>' first")
	}
	
	fmt.Printf("Running introspection on %s...\n", c.session.Target)
	
	// Simulate introspection results
	c.session.Schema = map[string]interface{}{
		"types": []string{"User", "Post", "Comment"},
		"queries": []string{"user", "users", "post", "posts"},
		"mutations": []string{"createUser", "updateUser", "deleteUser"},
	}
	
	c.session.Queries = []string{"user(id: ID!)", "users(limit: Int)", "post(id: ID!)", "posts(authorId: ID)"}
	c.session.Mutations = []string{"createUser(input: UserInput!)", "updateUser(id: ID!, input: UserInput!)", "deleteUser(id: ID!)"}
	
	// Add discovery finding
	finding := intel.Finding{
		Type:        "schema_discovery",
		Severity:    "info",
		Title:       "GraphQL Schema Discovered",
		Description: "Successfully retrieved GraphQL schema via introspection",
		Location:    c.session.Target,
		Evidence: map[string]interface{}{
			"types_count":     3,
			"queries_count":   4,
			"mutations_count": 3,
		},
		Timestamp: time.Now(),
	}
	c.session.Discoveries = append(c.session.Discoveries, finding)
	
	fmt.Printf("✓ Schema discovered!\n")
	fmt.Printf("  • Types: %d\n", 3)
	fmt.Printf("  • Queries: %d\n", 4)
	fmt.Printf("  • Mutations: %d\n", 3)
	fmt.Printf("\nTry: %sintel analyze%s for AI insights\n", output.YellowColor, output.Reset)
	
	return nil
}

func (c *IntrospectCommand) Description() string {
	return "Discover GraphQL schema via introspection"
}

// Other command implementations...
type QueryCommand struct{ session *GraphQLSession }
func (c *QueryCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: query <graphql_query>")
	}
	query := strings.Join(args, " ")
	fmt.Printf("Executing query: %s\n", query)
	fmt.Printf("✓ Query executed successfully (mock)\n")
	return nil
}
func (c *QueryCommand) Description() string { return "Execute GraphQL queries" }

type ScanCommand struct{ session *GraphQLSession }
func (c *ScanCommand) Execute(args []string) error {
	fmt.Printf("Running automated security scans...\n")
	
	// Simulate finding a vulnerability
	finding := intel.Finding{
		Type:        "vulnerability",
		Severity:    "high",
		Title:       "Authorization Bypass in User Query",
		Description: "The user query does not properly validate user permissions, allowing unauthorized access to user data",
		Location:    "Query.user",
		Evidence: map[string]interface{}{
			"query": "user(id: \"other_user_id\")",
			"response": "Returned data for unauthorized user",
		},
		Timestamp: time.Now(),
	}
	c.session.Discoveries = append(c.session.Discoveries, finding)
	
	fmt.Printf("%s! High severity vulnerability found%s\n", output.RedColor, output.Reset)
	fmt.Printf("Try: %sintel explain authorization bypass%s\n", output.YellowColor, output.Reset)
	return nil
}
func (c *ScanCommand) Description() string { return "Run automated GraphQL security scans" }

type ShowCommand struct{ session *GraphQLSession; state *config.State }
func (c *ShowCommand) Execute(args []string) error {
	fmt.Printf("\n%sSession Status%s\n", output.BoldColor, output.Reset)
	fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 14), output.Reset)
	fmt.Printf("Target: %s\n", c.session.Target)
	fmt.Printf("Schema: %s\n", map[bool]string{true: "✅ Discovered", false: "❌ Not discovered"}[len(c.session.Schema) > 0])
	fmt.Printf("Findings: %d\n", len(c.session.Discoveries))
	return nil
}
func (c *ShowCommand) Description() string { return "Display session information" }

type AuthCommand struct{ session *GraphQLSession }
func (c *AuthCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: auth <token>")
	}
	c.session.Token = args[0]
	c.session.Authenticated = true
	fmt.Printf("✓ Authentication token set\n")
	return nil
}
func (c *AuthCommand) Description() string { return "Set authentication token" }