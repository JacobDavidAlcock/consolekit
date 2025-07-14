package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// LoadWordlist loads a wordlist from a file, returning a slice of words
// Extracted and adapted from firescan's loadWordlist function
func LoadWordlist(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" { // Skip empty lines
			words = append(words, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return words, nil
}

// SaveWordlist saves a wordlist to a file
func SaveWordlist(filePath string, words []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, word := range words {
		if _, err := writer.WriteString(word + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return home, nil
}

// GetConfigDir returns a configuration directory path
func GetConfigDir(appName string) (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(home, ".config", appName)
	return configDir, nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath(appName string) (string, error) {
	configDir, err := GetConfigDir(appName)
	if err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "config.yaml"), nil
}

// ReadLines reads all lines from a file
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteLines writes lines to a file
func WriteLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// AppendToFile appends text to a file
func AppendToFile(filePath, text string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}