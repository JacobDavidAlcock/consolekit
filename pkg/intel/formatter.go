package intel

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// StreamingFormatter handles real-time markdown formatting and streaming
type StreamingFormatter struct {
	buffer       strings.Builder
	inCodeBlock  bool
	lastOutput   time.Time
	minDelay     time.Duration
}

// NewStreamingFormatter creates a new streaming formatter
func NewStreamingFormatter() *StreamingFormatter {
	return &StreamingFormatter{
		minDelay: 30 * time.Millisecond, // Minimum delay between characters for typing effect
	}
}

// ProcessToken processes a single token from the LLM stream
func (f *StreamingFormatter) ProcessToken(token string) {
	f.buffer.WriteString(token)
	
	// For now, just show the raw text with typing effect
	// We'll format everything at the end
	f.showTypingEffect(token)
}

// showTypingEffect displays text with a typing animation
func (f *StreamingFormatter) showTypingEffect(text string) {
	if time.Since(f.lastOutput) < f.minDelay {
		time.Sleep(f.minDelay - time.Since(f.lastOutput))
	}
	
	for _, char := range text {
		fmt.Print(string(char))
		f.lastOutput = time.Now()
		time.Sleep(15 * time.Millisecond) // Typing effect delay
	}
}

// formatLine applies markdown formatting to a line
func (f *StreamingFormatter) formatLine(line string) string {
	// Handle code blocks
	if strings.HasPrefix(strings.TrimSpace(line), "```") {
		f.inCodeBlock = !f.inCodeBlock
		if f.inCodeBlock {
			return fmt.Sprintf("\n%s╭─ Code Block ──────────────────────────────────────────────────╮%s", 
				output.CyanColor, output.Reset)
		} else {
			return fmt.Sprintf("%s╰──────────────────────────────────────────────────────────────╯%s", 
				output.CyanColor, output.Reset)
		}
	}
	
	if f.inCodeBlock {
		return fmt.Sprintf("%s│ %s%s", output.CyanColor, line, output.Reset)
	}
	
	// Handle headers with ASCII art
	if strings.HasPrefix(strings.TrimSpace(line), "## ") {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "## "))
		length := len(text)
		if length > 60 {
			length = 60
		}
		border := strings.Repeat("─", length+4)
		return fmt.Sprintf("\n%s╭─%s─╮%s\n%s│ %s%s%s │%s\n%s╰─%s─╯%s", 
			output.CyanColor, border, output.Reset,
			output.CyanColor, output.BoldColor, text, output.Reset, output.CyanColor,
			output.CyanColor, border, output.Reset)
	}
	
	if strings.HasPrefix(strings.TrimSpace(line), "# ") {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "# "))
		length := len(text)
		if length > 60 {
			length = 60
		}
		border := strings.Repeat("═", length+4)
		return fmt.Sprintf("\n%s╔═%s═╗%s\n%s║ %s%s%s ║%s\n%s╚═%s═╝%s", 
			output.CyanColor, border, output.Reset,
			output.CyanColor, output.BoldColor, text, output.Reset, output.CyanColor,
			output.CyanColor, border, output.Reset)
	}
	
	// Handle bullet points with enhanced ASCII
	if strings.HasPrefix(strings.TrimSpace(line), "- ") {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "- "))
		return fmt.Sprintf("  %s▸%s %s", output.YellowColor, output.Reset, f.formatInlineMarkdown(text))
	}
	
	if strings.HasPrefix(strings.TrimSpace(line), "* ") {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "* "))
		return fmt.Sprintf("  %s▸%s %s", output.YellowColor, output.Reset, f.formatInlineMarkdown(text))
	}
	
	// Handle numbered lists with enhanced formatting
	if matched, _ := regexp.MatchString(`^\s*\d+\.\s`, line); matched {
		re := regexp.MustCompile(`^(\s*)(\d+\.)(\s*)(.*)$`)
		matches := re.FindStringSubmatch(line)
		if len(matches) == 5 {
			indent := matches[1]
			number := matches[2]
			space := matches[3]
			text := matches[4]
			return fmt.Sprintf("%s%s%s%s%s %s", 
				indent, output.GreenColor, number, output.Reset, space, f.formatInlineMarkdown(text))
		}
	}
	
	// Handle long lines by wrapping them
	if len(line) > 100 {
		return f.wrapLongLine(line)
	}
	
	// Regular line with inline formatting
	return f.formatInlineMarkdown(line)
}

// formatInlineMarkdown handles inline markdown formatting
func (f *StreamingFormatter) formatInlineMarkdown(text string) string {
	// Clean up the text first - remove common LLM artifacts
	text = f.cleanLLMOutput(text)
	
	// Handle bold text **text**
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = boldRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "*")
		return fmt.Sprintf("%s%s%s", output.BoldColor, content, output.Reset)
	})
	
	// Handle italic text *text* (but avoid conflicting with bold)
	italicRegex := regexp.MustCompile(`(?:\*\*)?\*([^*]+)\*(?:\*\*)?`)
	text = italicRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Skip if this is part of a bold pattern
		if strings.Contains(match, "**") {
			return match
		}
		content := strings.Trim(match, "*")
		return fmt.Sprintf("%s%s%s", output.CyanColor, content, output.Reset)
	})
	
	// Handle inline code `code`
	codeRegex := regexp.MustCompile("`([^`]+)`")
	text = codeRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "`")
		return fmt.Sprintf("%s%s%s", output.YellowColor, content, output.Reset)
	})
	
	return text
}

// cleanLLMOutput removes common artifacts and symbols from LLM output
func (f *StreamingFormatter) cleanLLMOutput(text string) string {
	// Remove common LLM artifacts and formatting symbols
	patterns := []struct {
		pattern string
		replacement string
	}{
		// Remove markdown table separators that appear randomly
		{`\|\s*-+\s*\|`, ""},
		{`\|\s*:?-+:?\s*\|`, ""},
		
		// Clean up excessive formatting
		{`\*{3,}`, "**"}, // Reduce multiple asterisks to bold
		{`_{3,}`, "__"},  // Reduce multiple underscores
		
		// Remove orphaned formatting characters
		{`(?:^|\s)\*(?:\s|$)`, " "},     // Standalone asterisks
		{`(?:^|\s)_(?:\s|$)`, " "},      // Standalone underscores
		{`(?:^|\s)\|(?:\s|$)`, " "},     // Standalone pipes
		
		// Clean up extra whitespace
		{`\s{3,}`, "  "},                // Reduce multiple spaces
		{`\n{3,}`, "\n\n"},              // Reduce multiple newlines
		
		// Remove common formatting artifacts
		{`^\s*[\|\-\+\=]{1,}\s*$`, ""}, // Lines with only formatting chars
		{`^\s*\.\.\.\s*$`, ""},          // Lines with just dots
		
		// Fix broken formatting
		{`\*\s+\*`, ""},                 // Broken asterisks with spaces
		{`_\s+_`, ""},                   // Broken underscores with spaces
		
		// Remove verbose phrases common in LLM output
		{`(?i)here's?\s+(?:a\s+)?(?:comprehensive\s+)?(?:summary|overview|breakdown|explanation)\s*:?\s*`, ""},
		{`(?i)let me\s+(?:provide|explain|show|give)\s+(?:you\s+)?(?:a\s+)?`, ""},
		{`(?i)to\s+(?:help\s+)?(?:you\s+)?(?:understand|get\s+started|begin)`, ""},
		{`(?i)(?:as\s+)?(?:you\s+)?(?:can\s+)?(?:see|notice|observe)`, ""},
		{`(?i)(?:it's\s+)?(?:important\s+)?(?:to\s+)?(?:note|remember|keep\s+in\s+mind)\s+that`, ""},
		{`(?i)(?:please\s+)?(?:also\s+)?(?:note\s+)?(?:that\s+)?(?:these\s+)?(?:steps\s+)?(?:are\s+)?(?:designed\s+)?(?:to\s+)?`, ""},
		
		// Remove redundant transitions
		{`(?i)(?:now\s+)?(?:let's\s+)?(?:next\s+)?(?:step\s+)?(?:we'll\s+)?(?:move\s+)?(?:to\s+)?(?:the\s+)?(?:next\s+)?(?:part\s+)?`, ""},
		{`(?i)(?:in\s+)?(?:this\s+)?(?:section\s+)?(?:we\s+)?(?:will\s+)?(?:cover\s+)?(?:discuss\s+)?`, ""},
		
		// Clean up numbered items format
		{`(?i)(?:step\s+)?(\d+)\.?\s*[:.]?\s*`, "$1. "},
		
		// Remove filler words and phrases
		{`(?i)(?:essentially|basically|fundamentally|primarily|generally|typically|usually|often|commonly)`, ""},
		{`(?i)(?:as\s+mentioned|as\s+noted|as\s+discussed)(?:\s+(?:before|above|previously|earlier))?`, ""},
	}
	
	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		text = re.ReplaceAllString(text, p.replacement)
	}
	
	// Final cleanup - remove leading/trailing whitespace from lines
	lines := strings.Split(text, "\n")
	cleanedLines := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and lines with only formatting
		if line != "" && !regexp.MustCompile(`^[\s\-\*\|\.]+$`).MatchString(line) {
			cleanedLines = append(cleanedLines, line)
		}
	}
	
	text = strings.Join(cleanedLines, "\n")
	
	// Remove empty lines at start and end
	text = strings.TrimSpace(text)
	
	return text
}

// wrapLongLine wraps long lines for better readability
func (f *StreamingFormatter) wrapLongLine(line string) string {
	const maxWidth = 80
	words := strings.Fields(line)
	if len(words) == 0 {
		return line
	}
	
	var result strings.Builder
	var currentLine strings.Builder
	currentLength := 0
	
	for _, word := range words {
		wordLen := len(word)
		
		// If adding this word would exceed max width, start a new line
		if currentLength > 0 && currentLength+wordLen+1 > maxWidth {
			result.WriteString(f.formatInlineMarkdown(currentLine.String()))
			result.WriteString("\n")
			currentLine.Reset()
			currentLength = 0
		}
		
		// Add word to current line
		if currentLength > 0 {
			currentLine.WriteString(" ")
			currentLength++
		}
		currentLine.WriteString(word)
		currentLength += wordLen
	}
	
	// Add the last line
	if currentLength > 0 {
		result.WriteString(f.formatInlineMarkdown(currentLine.String()))
	}
	
	return result.String()
}

// Complete finishes the formatting process
func (f *StreamingFormatter) Complete() {
	// Format the complete text now
	completeText := f.buffer.String()
	formatted := f.formatCompleteText(completeText)
	
	// Clear the screen line and print formatted version
	fmt.Printf("\r%s\r", strings.Repeat(" ", 80))
	fmt.Print(formatted)
	fmt.Print("\n")
}

// FormatAndDisplayResponse formats and displays a complete response with proper line handling
func (f *StreamingFormatter) FormatAndDisplayResponse(content string) {
	// Clean the content first
	content = f.cleanLLMOutput(content)
	
	// Truncate if too long (Claude-like behavior)
	content = f.truncateResponse(content)
	
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		// Handle empty lines
		if strings.TrimSpace(line) == "" {
			fmt.Printf("\n")
			continue
		}
		
		// Format the line
		formatted := f.formatLine(line)
		
		// Display the formatted line
		if formatted != "" {
			fmt.Printf("%s\n", formatted)
		}
		
		// Small delay for readability
		time.Sleep(60 * time.Millisecond)
	}
}

// truncateResponse truncates responses that are too long for CLI usage
func (f *StreamingFormatter) truncateResponse(content string) string {
	// Count words to stay within reasonable limits
	words := strings.Fields(content)
	maxWords := 200 // Claude-like conciseness
	
	if len(words) <= maxWords {
		return content
	}
	
	// Truncate at word boundary
	truncated := strings.Join(words[:maxWords], " ")
	
	// Try to end at a natural break point
	lastPeriod := strings.LastIndex(truncated, ".")
	lastNewline := strings.LastIndex(truncated, "\n")
	
	cutPoint := lastPeriod
	if lastNewline > lastPeriod {
		cutPoint = lastNewline
	}
	
	// If we found a natural break point, use it
	if cutPoint > len(truncated)/2 {
		truncated = truncated[:cutPoint]
	}
	
	// Add truncation indicator
	truncated += "\n\n" + fmt.Sprintf("%s[Response truncated for readability]%s", 
		output.YellowColor, output.Reset)
	
	return truncated
}

// formatCompleteText formats the complete response text
func (f *StreamingFormatter) formatCompleteText(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder
	
	for _, line := range lines {
		formatted := f.formatLine(line)
		result.WriteString(formatted)
		if !strings.HasSuffix(formatted, "\n") {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// GetPersonalityMessage returns a contextual personality message
func GetPersonalityMessage(context string) string {
	messages := map[string][]string{
		"thinking": {
			"Processing request...",
			"Analyzing situation...",
			"Connecting the dots...",
			"Diving deep into this...",
		},
		"analyzing": {
			"Analyzing session...",
			"Crunching data...",
			"Investigating details...",
			"Examining evidence...",
			"Reviewing patterns...",
		},
		"suggesting": {
			"Brainstorming ideas...",
			"Finding best approach...",
			"Mapping out next steps...",
			"Planning strategy...",
			"Curating recommendations...",
		},
		"explaining": {
			"Preparing explanation...",
			"Gathering knowledge...",
			"Consulting database...",
			"Breaking down concept...",
			"Crafting explanation...",
		},
		"initializing": {
			"Initializing Intel system...",
			"Starting AI services...",
			"Loading components...",
			"Coming online...",
			"Preparing for action...",
		},
		"downloading": {
			"Fetching model...",
			"Downloading from server...",
			"Unpacking capabilities...",
			"Synchronizing data...",
			"Pulling model files...",
		},
	}
	
	if msgs, exists := messages[context]; exists {
		// Use simple random selection based on time
		index := int(time.Now().UnixNano()) % len(msgs)
		return msgs[index]
	}
	
	return "Working on it..."
}

// ShowPersonalityMessage displays a personality message with animation
func ShowPersonalityMessage(context string) {
	message := GetPersonalityMessage(context)
	fmt.Printf("%s%s%s", output.CyanColor, message, output.Reset)
	
	// More sophisticated animation like Claude CLI
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	for i := 0; i < 20; i++ {
		fmt.Printf("\r%s%s%s %s", output.CyanColor, message, output.Reset, spinner[i%len(spinner)])
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Print("\r")
}

// DownloadTracker tracks and displays download progress
type DownloadTracker struct {
	totalSize    int64
	downloaded   int64
	startTime    time.Time
	lastUpdate   time.Time
	lastPrint    time.Time
	currentPhase string
	speedSamples []float64
	maxSamples   int
}

// NewDownloadTracker creates a new download tracker
func NewDownloadTracker() *DownloadTracker {
	return &DownloadTracker{
		startTime:    time.Now(),
		lastUpdate:   time.Now(),
		lastPrint:    time.Now(),
		speedSamples: make([]float64, 0),
		maxSamples:   10,
	}
}

// Update processes a progress response from Ollama
func (d *DownloadTracker) Update(resp interface{}) {
	// Try to extract progress information safely
	// Since the API might change, we'll use interface{} and type assertions
	now := time.Now()
	
	// Only update display every 100ms to avoid spam
	if now.Sub(d.lastPrint) < 100*time.Millisecond {
		return
	}
	d.lastPrint = now
	
	// Try to extract status and progress information
	if respMap, ok := resp.(map[string]interface{}); ok {
		if status, exists := respMap["status"].(string); exists {
			d.currentPhase = status
		}
		
		if completed, exists := respMap["completed"].(int64); exists {
			d.downloaded = completed
		}
		
		if total, exists := respMap["total"].(int64); exists {
			d.totalSize = total
		}
	}
	
	// Calculate speed
	if d.totalSize > 0 && d.downloaded > 0 {
		speed := d.calculateSpeed()
		percentage := float64(d.downloaded) / float64(d.totalSize) * 100
		eta := d.calculateETA(speed)
		
		// Clear line and show progress
		fmt.Printf("\r%s[%s] %.1f%% (%.1f MB/s) ETA: %s%s", 
			output.CyanColor,
			d.createProgressBar(percentage),
			percentage,
			speed,
			eta,
			output.Reset)
	} else {
		// Show phase information when no progress data available
		fmt.Printf("\r%s%s...%s", output.CyanColor, d.currentPhase, output.Reset)
	}
	d.lastUpdate = now
}

// calculateSpeed calculates download speed in MB/s
func (d *DownloadTracker) calculateSpeed() float64 {
	if d.downloaded == 0 {
		return 0
	}
	
	elapsed := time.Since(d.startTime).Seconds()
	if elapsed == 0 {
		return 0
	}
	
	speed := float64(d.downloaded) / (1024 * 1024) / elapsed // MB/s
	
	// Add to samples for smoothing
	d.speedSamples = append(d.speedSamples, speed)
	if len(d.speedSamples) > d.maxSamples {
		d.speedSamples = d.speedSamples[1:]
	}
	
	// Return average speed
	sum := 0.0
	for _, s := range d.speedSamples {
		sum += s
	}
	return sum / float64(len(d.speedSamples))
}

// calculateETA calculates estimated time of arrival
func (d *DownloadTracker) calculateETA(speed float64) string {
	if speed == 0 || d.totalSize == 0 {
		return "calculating..."
	}
	
	remaining := d.totalSize - d.downloaded
	remainingMB := float64(remaining) / (1024 * 1024)
	eta := time.Duration(remainingMB/speed) * time.Second
	
	if eta > time.Hour {
		return fmt.Sprintf("%.1fh", eta.Hours())
	} else if eta > time.Minute {
		return fmt.Sprintf("%.1fm", eta.Minutes())
	} else {
		return fmt.Sprintf("%.0fs", eta.Seconds())
	}
}

// createProgressBar creates a visual progress bar
func (d *DownloadTracker) createProgressBar(percentage float64) string {
	width := 20
	filled := int(percentage / 100 * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

// Complete finishes the download tracking
func (d *DownloadTracker) Complete() {
	fmt.Printf("\r%80s\r", "") // Clear line
	elapsed := time.Since(d.startTime)
	fmt.Printf("%s✅ Download completed in %s%s\n", 
		output.GreenColor, elapsed.Round(time.Second), output.Reset)
}