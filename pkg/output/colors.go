package output

// ANSI color constants extracted from firescan
const (
	Reset       = "\033[0m"
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	CyanColor   = "\033[36m"
	BoldColor   = "\033[1m"
)

// Colorize wraps text with the specified color
func Colorize(text, color string) string {
	return color + text + Reset
}

// Red text
func Red(text string) string {
	return Colorize(text, RedColor)
}

// Green text
func Green(text string) string {
	return Colorize(text, GreenColor)
}

// Yellow text
func Yellow(text string) string {
	return Colorize(text, YellowColor)
}

// Cyan text
func Cyan(text string) string {
	return Colorize(text, CyanColor)
}

// Bold text
func Bold(text string) string {
	return Colorize(text, BoldColor)
}