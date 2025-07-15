package intel

import (
	"strings"
	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// StyleConstants provides consistent styling patterns for Intel CLI output
type StyleConstants struct{}

// GetStyleConstants returns the styling constants
func GetStyleConstants() *StyleConstants {
	return &StyleConstants{}
}

// Status indicators - minimal and professional
const (
	StatusSuccess  = "✓"
	StatusError    = "❌"
	StatusWarning  = "⚠️"
	StatusInfo     = "ℹ️"
	StatusProgress = "…"
)

// Professional ASCII art patterns
func (s *StyleConstants) CreateHeader(title string, style string) string {
	length := len(title)
	if length > 60 {
		length = 60
		title = title[:60]
	}
	
	switch style {
	case "main":
		// Main header with double lines
		border := strings.Repeat("═", length+4)
		return output.CyanColor + "╔═" + border + "═╗\n" +
			"║ " + output.BoldColor + title + output.Reset + output.CyanColor + " ║\n" +
			"╚═" + border + "═╝" + output.Reset
	case "section":
		// Section header with single lines
		border := strings.Repeat("─", length+4)
		return output.CyanColor + "╭─" + border + "─╮\n" +
			"│ " + output.BoldColor + title + output.Reset + output.CyanColor + " │\n" +
			"╰─" + border + "─╯" + output.Reset
	case "simple":
		// Simple underlined header
		underline := strings.Repeat("─", length)
		return output.BoldColor + title + output.Reset + "\n" +
			output.CyanColor + underline + output.Reset
	default:
		return output.BoldColor + title + output.Reset
	}
}

// Create consistent separators
func (s *StyleConstants) CreateSeparator(width int, style string) string {
	if width <= 0 {
		width = 50
	}
	
	switch style {
	case "double":
		return output.CyanColor + strings.Repeat("═", width) + output.Reset
	case "single":
		return output.CyanColor + strings.Repeat("─", width) + output.Reset
	case "dotted":
		return output.CyanColor + strings.Repeat("·", width) + output.Reset
	default:
		return output.CyanColor + strings.Repeat("─", width) + output.Reset
	}
}

// Format status messages consistently
func (s *StyleConstants) FormatStatus(message string, status string) string {
	var indicator string
	var color string
	
	switch status {
	case "success":
		indicator = StatusSuccess
		color = output.GreenColor
	case "error":
		indicator = StatusError
		color = output.RedColor
	case "warning":
		indicator = StatusWarning
		color = output.YellowColor
	case "info":
		indicator = StatusInfo
		color = output.CyanColor
	case "progress":
		indicator = StatusProgress
		color = output.CyanColor
	default:
		indicator = StatusInfo
		color = output.CyanColor
	}
	
	return color + indicator + " " + message + output.Reset
}

// Format bullet points consistently
func (s *StyleConstants) FormatBullet(text string) string {
	return output.YellowColor + "▸" + output.Reset + " " + text
}

// Format numbered items consistently
func (s *StyleConstants) FormatNumbered(number int, text string) string {
	return output.GreenColor + string(rune('0'+number)) + "." + output.Reset + " " + text
}

// Format code blocks consistently
func (s *StyleConstants) FormatCodeBlock(code string) string {
	lines := strings.Split(code, "\n")
	var result strings.Builder
	
	result.WriteString(output.CyanColor + "╭─ Code Block " + strings.Repeat("─", 50) + "╮\n")
	for _, line := range lines {
		result.WriteString("│ " + line + "\n")
	}
	result.WriteString("╰" + strings.Repeat("─", 63) + "╯" + output.Reset)
	
	return result.String()
}

// Format inline code consistently
func (s *StyleConstants) FormatInlineCode(code string) string {
	return output.YellowColor + "`" + code + "`" + output.Reset
}

// Format emphasis consistently
func (s *StyleConstants) FormatEmphasis(text string, level string) string {
	switch level {
	case "strong":
		return output.BoldColor + text + output.Reset
	case "emphasis":
		return output.CyanColor + text + output.Reset
	case "highlight":
		return output.YellowColor + text + output.Reset
	default:
		return text
	}
}

// Professional styling guidelines as comments for reference
/*
CLI STYLING GUIDELINES FOR INTEL SYSTEM:

1. STATUS INDICATORS:
   - Use minimal Unicode symbols: ✓ ❌ ⚠️ ℹ️ …
   - Avoid excessive emojis or decorative symbols
   - Keep status messages concise and informative

2. HEADERS AND STRUCTURE:
   - Use box-drawing characters for professional appearance
   - Main headers: double lines (═)
   - Section headers: single lines (─)
   - Simple headers: underlined text

3. COLORS:
   - Success: Green
   - Error: Red
   - Warning: Yellow
   - Info/Headers: Cyan
   - Emphasis: Bold
   - Code: Yellow background

4. CONTENT FORMATTING:
   - Bullet points: ▸ (professional arrow)
   - Numbered lists: colored numbers
   - Code blocks: bordered boxes
   - Inline code: backticks with color

5. SPACING AND LAYOUT:
   - Consistent indentation (2 spaces)
   - Proper line spacing between sections
   - Maximum line width: 80 characters
   - Use whitespace effectively for readability

6. PROFESSIONAL TONE:
   - Minimal emoji usage
   - Clear, concise messaging
   - Structured information hierarchy
   - Consistent terminology

7. INTERACTIVE ELEMENTS:
   - Tab completion for all commands
   - Clear command syntax in help
   - Contextual suggestions
   - Progress indicators for long operations
*/