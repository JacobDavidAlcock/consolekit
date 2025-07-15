package intel

import (
	"fmt"
	"strings"

	"github.com/jacobdavidalcock/consolekit/pkg/console"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// RegisterIntelCommands adds Intel commands to a ConsoleKit application
func RegisterIntelCommands(app *console.Console, intel *IntelSystem) {
	// Main intel command with subcommands
	app.AddCommand("intel", &IntelCommand{system: intel}, "AI-powered analysis and assistance")
}

// IntelCommand handles all intel subcommands
type IntelCommand struct {
	system *IntelSystem
}

// Execute handles intel command execution with subcommands
func (c *IntelCommand) Execute(args []string) error {
	if len(args) == 0 {
		c.showHelp()
		return nil
	}

	subcommand := strings.ToLower(args[0])
	subArgs := args[1:]

	switch subcommand {
	case "start", "init":
		return c.handleStart(subArgs)
	case "analyze", "analysis":
		return c.handleAnalyze(subArgs)
	case "suggest", "suggestions":
		return c.handleSuggest(subArgs)
	case "explain", "explanation":
		return c.handleExplain(subArgs)
	case "status":
		return c.handleStatus(subArgs)
	case "context":
		return c.handleContext(subArgs)
	case "validate":
		return c.handleValidate(subArgs)
	case "help":
		if len(subArgs) > 0 && subArgs[0] == "errors" {
			ShowQuickHelp()
		} else {
			c.showHelp()
		}
		return nil
	default:
		return fmt.Errorf("unknown intel subcommand: %s. Use 'intel help' for available commands", subcommand)
	}
}

// Description returns the command description
func (c *IntelCommand) Description() string {
	return "AI-powered analysis and assistance"
}

// handleStart initializes the Intel system
func (c *IntelCommand) handleStart(args []string) error {
	// Show personality message
	ShowPersonalityMessage("initializing")
	
	if err := c.system.Initialize(); err != nil {
		fmt.Printf("\r%80s\r", "") // Clear personality message
		if intelErr, ok := err.(*IntelError); ok {
			intelErr.Display()
		} else {
			style := GetStyleConstants()
			fmt.Printf("%s\n", style.FormatStatus("Failed to initialize Intel system: "+err.Error(), "error"))
			fmt.Printf("%s\n", style.FormatStatus("Try: `ollama serve` or check your configuration", "info"))
		}
		return err
	}

	style := GetStyleConstants()
	fmt.Printf("\r%s\n", style.FormatStatus("Intel AI system initialized successfully", "success"))
	
	// Show available providers
	if len(c.system.providers) > 0 {
		fmt.Printf("\nActive context providers:\n")
		for _, provider := range c.system.providers {
			fmt.Printf("  %s\n", style.FormatBullet(provider.Name()))
		}
	}

	fmt.Printf("\n%sTry: intel analyze%s\n", output.CyanColor, output.Reset)
	return nil
}

// handleAnalyze performs AI analysis of the current session
func (c *IntelCommand) handleAnalyze(args []string) error {
	if !c.system.IsInitialized() {
		return fmt.Errorf("Intel system not initialized. Run 'intel start' first")
	}

	userPrompt := "Analyze the current session"
	if len(args) > 0 {
		userPrompt = strings.Join(args, " ")
	}

	// Show personality message
	ShowPersonalityMessage("analyzing")
	
	fmt.Printf("\n%sIntel Analysis%s\n", output.BoldColor, output.Reset)
	fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 15), output.Reset)
	
	// Use streaming analysis
	if err := c.system.AnalyzeWithStreaming(userPrompt); err != nil {
		if intelErr, ok := err.(*IntelError); ok {
			intelErr.Display()
		} else {
			fmt.Printf("%s❌ Analysis failed: %s%s\n", 
				output.RedColor, err.Error(), output.Reset)
		}
		return err
	}
	
	return nil
}

// handleSuggest provides AI-generated suggestions
func (c *IntelCommand) handleSuggest(args []string) error {
	if !c.system.IsInitialized() {
		return fmt.Errorf("Intel system not initialized. Run 'intel start' first")
	}

	context := "current session"
	if len(args) > 0 {
		context = strings.Join(args, " ")
	}

	// Show personality message
	ShowPersonalityMessage("suggesting")
	
	fmt.Printf("\n%sIntel Suggestions%s\n", output.BoldColor, output.Reset)
	fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 17), output.Reset)
	
	// Use streaming suggestions
	if err := c.system.SuggestWithStreaming(context); err != nil {
		if intelErr, ok := err.(*IntelError); ok {
			intelErr.Display()
		} else {
			fmt.Printf("%s❌ Suggestion generation failed: %s%s\n", 
				output.RedColor, err.Error(), output.Reset)
		}
		return err
	}
	
	return nil
}

// handleExplain provides detailed explanations
func (c *IntelCommand) handleExplain(args []string) error {
	if !c.system.IsInitialized() {
		return fmt.Errorf("Intel system not initialized. Run 'intel start' first")
	}

	if len(args) == 0 {
		return fmt.Errorf("please specify what you'd like explained. Usage: intel explain <topic>")
	}

	topic := strings.Join(args, " ")
	
	// Show personality message
	ShowPersonalityMessage("explaining")
	
	fmt.Printf("\n%sIntel Explanation%s\n", output.BoldColor, output.Reset)
	fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 17), output.Reset)
	fmt.Printf("%sTopic: %s%s\n\n", output.YellowColor, topic, output.Reset)
	
	// Use streaming explanation
	if err := c.system.ExplainWithStreaming(topic); err != nil {
		if intelErr, ok := err.(*IntelError); ok {
			intelErr.Display()
		} else {
			fmt.Printf("%s❌ Explanation failed: %s%s\n", 
				output.RedColor, err.Error(), output.Reset)
		}
		return err
	}
	
	return nil
}

// handleStatus shows Intel system status
func (c *IntelCommand) handleStatus(args []string) error {
	fmt.Printf("\n%s🤖 Intel System Status:%s\n", output.BoldColor, output.Reset)
	
	// Show Ollama status
	ollamaStatus, err := c.system.GetOllamaStatus()
	if err != nil {
		fmt.Printf("Ollama: %s%s%s\n", output.RedColor, ollamaStatus, output.Reset)
	} else {
		fmt.Printf("Ollama: %s\n", ollamaStatus)
	}
	
	// Show Intel system status
	if c.system.IsInitialized() {
		fmt.Printf("Intel: %s✅ Active%s\n", output.GreenColor, output.Reset)
		fmt.Printf("Model: %s%s%s\n", output.CyanColor, c.system.config.Model, output.Reset)
		fmt.Printf("URL: %s%s%s\n", output.CyanColor, c.system.config.OllamaURL, output.Reset)
	} else {
		fmt.Printf("Intel: %s❌ Not initialized%s\n", output.RedColor, output.Reset)
		fmt.Printf("Run '%sintel start%s' to initialize\n", output.YellowColor, output.Reset)
	}
	
	fmt.Printf("Providers: %s%d registered%s\n", output.CyanColor, len(c.system.providers), output.Reset)
	for _, provider := range c.system.providers {
		fmt.Printf("  • %s%s%s\n", output.YellowColor, provider.Name(), output.Reset)
	}
	
	// Show recent actions count
	c.system.context.mu.RLock()
	actionCount := len(c.system.context.RecentActions)
	c.system.context.mu.RUnlock()
	
	fmt.Printf("Context: %s%d recent actions%s\n", output.CyanColor, actionCount, output.Reset)
	fmt.Printf("Session: %s%s%s\n", output.CyanColor, c.system.context.StartTime.Format("15:04:05"), output.Reset)
	
	// Show context manager stats
	if c.system.IsInitialized() {
		stats := c.system.GetContextStats()
		fmt.Printf("\n%sContext Manager:%s\n", output.BoldColor, output.Reset)
		fmt.Printf("Tokens: %s%d/%d (%.1f%%)%s\n", 
			output.CyanColor, 
			stats["current_tokens"], 
			stats["max_tokens"], 
			stats["utilization"].(float64)*100, 
			output.Reset)
		fmt.Printf("Items: %s%d total%s\n", 
			output.CyanColor, 
			stats["total_items"], 
			output.Reset)
	}
	
	return nil
}

// showHelp displays help for intel commands
func (c *IntelCommand) showHelp() {
	style := GetStyleConstants()
	
	fmt.Printf("\n%s\n", style.CreateHeader("Intel AI Assistant Commands", "main"))
	fmt.Printf("  %sstart%s            Initialize the Intel AI system\n", output.GreenColor, output.Reset)
	fmt.Printf("  %sanalyze [query]%s   Analyze current session or specific query\n", output.GreenColor, output.Reset)
	fmt.Printf("  %ssuggest [context]%s Get AI suggestions for next steps\n", output.GreenColor, output.Reset)
	fmt.Printf("  %sexplain <topic>%s   Get detailed explanation of a concept\n", output.GreenColor, output.Reset)
	fmt.Printf("  %sstatus%s           Show Intel system status and configuration\n", output.GreenColor, output.Reset)
	fmt.Printf("  %scontext%s          Manage context (clear, stats, limit)\n", output.GreenColor, output.Reset)
	fmt.Printf("  %svalidate%s         Validate configuration (model, url, rules)\n", output.GreenColor, output.Reset)
	fmt.Printf("  %shelp%s             Show this help message\n", output.GreenColor, output.Reset)
	
	fmt.Printf("\n%s\n", style.CreateHeader("Examples", "section"))
	fmt.Printf("  intel start\n")
	fmt.Printf("  intel analyze\n")
	fmt.Printf("  intel suggest next steps\n")
	fmt.Printf("  intel explain GraphQL injection\n")
	fmt.Printf("  intel status\n")
	
	fmt.Printf("\n%s\n", style.FormatStatus("Intel requires Ollama to be running at " + c.system.config.OllamaURL, "info"))
}

// handleContext manages context information
func (c *IntelCommand) handleContext(args []string) error {
	if !c.system.IsInitialized() {
		return fmt.Errorf("Intel system not initialized. Run 'intel start' first")
	}

	if len(args) == 0 {
		// Show context summary
		fmt.Printf("\n%sContext Summary%s\n", output.BoldColor, output.Reset)
		fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 15), output.Reset)
		fmt.Printf("%s\n", c.system.GetContextSummary())
		return nil
	}

	subcommand := strings.ToLower(args[0])
	switch subcommand {
	case "clear":
		c.system.ClearContext()
		fmt.Printf("%s✓ Context cleared%s\n", output.GreenColor, output.Reset)
	case "stats":
		stats := c.system.GetContextStats()
		fmt.Printf("\n%sContext Statistics%s\n", output.BoldColor, output.Reset)
		fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 18), output.Reset)
		fmt.Printf("Total items: %d\n", stats["total_items"])
		fmt.Printf("Current tokens: %d/%d (%.1f%%)\n", 
			stats["current_tokens"], 
			stats["max_tokens"], 
			stats["utilization"].(float64)*100)
		
		if byType, ok := stats["by_type"].(map[string]int); ok {
			fmt.Printf("\nBy type:\n")
			for typeName, count := range byType {
				fmt.Printf("  • %s: %d\n", typeName, count)
			}
		}
	case "limit":
		if len(args) < 2 {
			return fmt.Errorf("usage: intel context limit <tokens>")
		}
		
		var limit int
		if _, err := fmt.Sscanf(args[1], "%d", &limit); err != nil {
			return fmt.Errorf("invalid token limit: %s", args[1])
		}
		
		if limit < 1000 || limit > 32000 {
			return fmt.Errorf("token limit must be between 1000 and 32000")
		}
		
		c.system.SetMaxTokens(limit)
		fmt.Printf("%s✓ Token limit set to %d%s\n", output.GreenColor, limit, output.Reset)
	default:
		return fmt.Errorf("unknown context subcommand: %s. Use 'clear', 'stats', or 'limit'", subcommand)
	}
	
	return nil
}

// handleValidate validates configuration
func (c *IntelCommand) handleValidate(args []string) error {
	validator := NewConfigValidator()
	
	if len(args) == 0 {
		// Validate current configuration
		fmt.Printf("\n%sValidating Intel Configuration%s\n", output.BoldColor, output.Reset)
		fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("=", 30), output.Reset)
		
		if err := validator.ValidateConfig(c.system.config); err != nil {
			if intelErr, ok := err.(*IntelError); ok {
				intelErr.Display()
			} else {
				fmt.Printf("%s❌ Validation failed: %s%s\n", 
					output.RedColor, err.Error(), output.Reset)
			}
			return err
		}
		
		fmt.Printf("%s✓ Configuration is valid%s\n", output.GreenColor, output.Reset)
		return nil
	}
	
	subcommand := strings.ToLower(args[0])
	switch subcommand {
	case "model":
		if len(args) < 2 {
			return fmt.Errorf("usage: intel validate model <model_name>")
		}
		
		modelName := args[1]
		fmt.Printf("\n%sValidating Model: %s%s\n", output.BoldColor, modelName, output.Reset)
		fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("-", 20), output.Reset)
		
		if err := validator.ValidateModel(modelName); err != nil {
			if intelErr, ok := err.(*IntelError); ok {
				intelErr.Display()
			} else {
				fmt.Printf("%s❌ Model validation failed: %s%s\n", 
					output.RedColor, err.Error(), output.Reset)
			}
			return err
		}
		
		// Check system requirements
		if err := validator.ValidateSystemRequirements(modelName); err != nil {
			if intelErr, ok := err.(*IntelError); ok {
				intelErr.Display()
			} else {
				fmt.Printf("%s⚠️  System requirements: %s%s\n", 
					output.YellowColor, err.Error(), output.Reset)
			}
		}
		
		fmt.Printf("%s✓ Model is valid%s\n", output.GreenColor, output.Reset)
		
		// Show model info
		if info, exists := validator.GetModelInfo(modelName); exists {
			fmt.Printf("\n%sModel Info:%s\n", output.BoldColor, output.Reset)
			fmt.Printf("Size: %s\n", info.Size)
			fmt.Printf("Description: %s\n", info.Description)
			fmt.Printf("Specialty: %s\n", info.Specialty)
			fmt.Printf("Min RAM: %dGB\n", info.MinRAM)
		}
		
	case "url":
		if len(args) < 2 {
			return fmt.Errorf("usage: intel validate url <url>")
		}
		
		url := args[1]
		fmt.Printf("\n%sValidating URL: %s%s\n", output.BoldColor, url, output.Reset)
		fmt.Printf("%s%s%s\n", output.CyanColor, strings.Repeat("-", 18), output.Reset)
		
		if err := validator.ValidateURL(url); err != nil {
			if intelErr, ok := err.(*IntelError); ok {
				intelErr.Display()
			} else {
				fmt.Printf("%s❌ URL validation failed: %s%s\n", 
					output.RedColor, err.Error(), output.Reset)
			}
			return err
		}
		
		fmt.Printf("%s✓ URL is valid%s\n", output.GreenColor, output.Reset)
		
	case "rules":
		fmt.Printf("\n%s", validator.GetValidationSummary())
		
	default:
		return fmt.Errorf("unknown validate subcommand: %s. Use 'model', 'url', or 'rules'", subcommand)
	}
	
	return nil
}