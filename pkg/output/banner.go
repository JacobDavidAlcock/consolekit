package output

import (
	"fmt"
	"strings"
)

// PrintBanner displays an ASCII art banner with optional color
func PrintBanner(banner, color string) {
	if color != "" {
		fmt.Println(color + banner + Reset)
	} else {
		fmt.Println(banner)
	}
}

// CreateSimpleBanner creates a simple text banner with borders
func CreateSimpleBanner(title, subtitle string) string {
	maxLen := len(title)
	if len(subtitle) > maxLen {
		maxLen = len(subtitle)
	}
	
	border := strings.Repeat("=", maxLen+4)
	
	banner := fmt.Sprintf("%s\n= %s =\n", border, CenterText(title, maxLen))
	if subtitle != "" {
		banner += fmt.Sprintf("= %s =\n", CenterText(subtitle, maxLen))
	}
	banner += border
	
	return banner
}

// CreateBoxBanner creates a banner with box drawing characters
func CreateBoxBanner(title, subtitle string) string {
	maxLen := len(title)
	if len(subtitle) > maxLen {
		maxLen = len(subtitle)
	}
	
	width := maxLen + 4
	top := "┌" + strings.Repeat("─", width-2) + "┐"
	bottom := "└" + strings.Repeat("─", width-2) + "┘"
	
	banner := top + "\n"
	banner += fmt.Sprintf("│ %s │\n", CenterText(title, width-4))
	if subtitle != "" {
		banner += fmt.Sprintf("│ %s │\n", CenterText(subtitle, width-4))
	}
	banner += bottom
	
	return banner
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	
	padding := width - len(text)
	leftPad := padding / 2
	rightPad := padding - leftPad
	
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// GenerateConsoleBanner creates a banner for console applications
func GenerateConsoleBanner(appName, description string) string {
	banner := CreateBoxBanner(appName, description)
	return CyanColor + banner + Reset
}

// PrintWelcome prints a welcome message with app info
func PrintWelcome(appName, version, description string) {
	banner := fmt.Sprintf(`
┌─────────────────────────────────────────┐
│ %s │
│ %s │
│ %s │
└─────────────────────────────────────────┘
`, 
		CenterText(BoldColor+appName+Reset, 39),
		CenterText("v"+version, 39),
		CenterText(description, 39))
	
	fmt.Println(CyanColor + banner + Reset)
}

// PrintSeparator prints a visual separator
func PrintSeparator(char string, length int, color string) {
	separator := strings.Repeat(char, length)
	if color != "" {
		fmt.Println(color + separator + Reset)
	} else {
		fmt.Println(separator)
	}
}

// PrintSection prints a section header
func PrintSection(title string) {
	fmt.Printf("\n%s=== %s ===%s\n", BoldColor, title, Reset)
}