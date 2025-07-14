package utils

import (
	"strings"
)

// MaskString hides the middle of a string for secure display
// Extracted from firescan's maskString function
func MaskString(s string, prefixLen, suffixLen int) string {
	if len(s) < prefixLen+suffixLen {
		return "..."
	}
	return s[:prefixLen] + "..." + s[len(s)-suffixLen:]
}

// GenerateCaseVariations takes a word and returns a slice with its lowercase, PascalCase, and UPPERCASE variations
// Extracted from firescan's generateCaseVariations function
func GenerateCaseVariations(word string) []string {
	if len(word) == 0 {
		return []string{}
	}
	
	variationsSet := make(map[string]bool)
	variationsSet[strings.ToLower(word)] = true
	variationsSet[strings.ToUpper(string(word[0]))+strings.ToLower(word[1:])] = true
	variationsSet[strings.ToUpper(word)] = true
	
	result := make([]string, 0, len(variationsSet))
	for v := range variationsSet {
		result = append(result, v)
	}
	return result
}

// Truncate truncates a string to a maximum length with optional suffix
func Truncate(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	
	if len(suffix) >= maxLen {
		return suffix[:maxLen]
	}
	
	return s[:maxLen-len(suffix)] + suffix
}

// PadRight pads a string to the right with spaces
func PadRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// PadLeft pads a string to the left with spaces
func PadLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

// Center centers a string within a given width
func Center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	
	padding := width - len(s)
	leftPad := padding / 2
	rightPad := padding - leftPad
	
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// RemoveEmpty removes empty strings from a slice
func RemoveEmpty(strings []string) []string {
	var result []string
	for _, s := range strings {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsIgnoreCase checks if a slice contains a string (case insensitive)
func ContainsIgnoreCase(slice []string, item string) bool {
	lower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == lower {
			return true
		}
	}
	return false
}