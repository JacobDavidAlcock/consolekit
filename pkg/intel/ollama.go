package intel

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jacobdavidalcock/consolekit/pkg/output"
)

// OllamaManager handles Ollama installation and service management
type OllamaManager struct {
	installDir   string
	binaryPath   string
	serviceURL   string
	downloadURL  string
	isInstalled  bool
	isRunning    bool
}

// NewOllamaManager creates a new Ollama manager
func NewOllamaManager() *OllamaManager {
	manager := &OllamaManager{
		serviceURL: "http://localhost:11434",
	}
	
	// Set platform-specific paths and URLs
	manager.setPlatformDefaults()
	
	return manager
}

// setPlatformDefaults sets OS-specific default paths and URLs
func (om *OllamaManager) setPlatformDefaults() {
	switch runtime.GOOS {
	case "windows":
		om.installDir = filepath.Join(os.Getenv("APPDATA"), "ollama")
		om.binaryPath = filepath.Join(om.installDir, "ollama.exe")
		om.downloadURL = "https://ollama.com/download/windows"
	case "darwin":
		om.installDir = "/usr/local/bin"
		om.binaryPath = filepath.Join(om.installDir, "ollama")
		om.downloadURL = "https://ollama.com/download/mac"
	case "linux":
		om.installDir = "/usr/local/bin"
		om.binaryPath = filepath.Join(om.installDir, "ollama")
		om.downloadURL = "https://ollama.com/download/linux"
	default:
		om.installDir = "/usr/local/bin"
		om.binaryPath = filepath.Join(om.installDir, "ollama")
		om.downloadURL = "https://ollama.com/download/linux"
	}
}

// EnsureOllamaAvailable checks if Ollama is installed and running, installs/starts if needed
func (om *OllamaManager) EnsureOllamaAvailable() error {
	// Check if Ollama is installed
	if err := om.checkInstallation(); err != nil {
		if !om.isInstalled {
			// Try to install Ollama
			if err := om.installOllama(); err != nil {
				return fmt.Errorf("failed to install Ollama: %w", err)
			}
		} else {
			return fmt.Errorf("Ollama installation check failed: %w", err)
		}
	}
	
	// Check if Ollama is running
	if err := om.checkService(); err != nil {
		if !om.isRunning {
			// Try to start Ollama
			if err := om.startOllama(); err != nil {
				return fmt.Errorf("failed to start Ollama: %w", err)
			}
		} else {
			return fmt.Errorf("Ollama service check failed: %w", err)
		}
	}
	
	return nil
}

// checkInstallation checks if Ollama is installed
func (om *OllamaManager) checkInstallation() error {
	// First check if ollama is in PATH
	if _, err := exec.LookPath("ollama"); err == nil {
		om.isInstalled = true
		om.binaryPath = "ollama" // Use PATH version
		return nil
	}
	
	// Check if it's in our expected location
	if _, err := os.Stat(om.binaryPath); err == nil {
		om.isInstalled = true
		return nil
	}
	
	om.isInstalled = false
	return fmt.Errorf("Ollama not found in PATH or %s", om.binaryPath)
}

// checkService checks if Ollama service is running
func (om *OllamaManager) checkService() error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(om.serviceURL + "/api/version")
	if err != nil {
		om.isRunning = false
		return fmt.Errorf("Ollama service not responding: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		om.isRunning = true
		return nil
	}
	
	om.isRunning = false
	return fmt.Errorf("Ollama service returned status: %d", resp.StatusCode)
}

// installOllama attempts to install Ollama automatically
func (om *OllamaManager) installOllama() error {
	fmt.Printf("%s🔧 Installing Ollama...%s\n", output.YellowColor, output.Reset)
	
	switch runtime.GOOS {
	case "linux":
		return om.installLinux()
	case "darwin":
		return om.installMacOS()
	case "windows":
		return om.installWindows()
	default:
		return fmt.Errorf("automatic installation not supported on %s", runtime.GOOS)
	}
}

// installLinux installs Ollama on Linux using the official script
func (om *OllamaManager) installLinux() error {
	// Use the official Ollama install script
	cmd := exec.Command("sh", "-c", "curl -fsSL https://ollama.com/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Ollama install script: %w", err)
	}
	
	// Update our paths
	om.binaryPath = "/usr/local/bin/ollama"
	om.isInstalled = true
	
	fmt.Printf("%s✅ Ollama installed successfully%s\n", output.GreenColor, output.Reset)
	return nil
}

// installMacOS installs Ollama on macOS
func (om *OllamaManager) installMacOS() error {
	// Check if Homebrew is available
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Printf("%s📦 Installing via Homebrew...%s\n", output.CyanColor, output.Reset)
		cmd := exec.Command("brew", "install", "ollama")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install via Homebrew: %w", err)
		}
		
		om.binaryPath = "/opt/homebrew/bin/ollama" // Common Homebrew path
		om.isInstalled = true
		
		fmt.Printf("%s✅ Ollama installed via Homebrew%s\n", output.GreenColor, output.Reset)
		return nil
	}
	
	// Fall back to manual download instructions
	return fmt.Errorf("Homebrew not found. Please install Ollama manually from: %s", om.downloadURL)
}

// installWindows provides Windows installation instructions
func (om *OllamaManager) installWindows() error {
	// Windows installation is more complex, provide instructions
	fmt.Printf("%s⚠️  Automatic installation not available on Windows%s\n", output.YellowColor, output.Reset)
	fmt.Printf("Please install Ollama manually:\n")
	fmt.Printf("1. Download from: %s%s%s\n", output.CyanColor, om.downloadURL, output.Reset)
	fmt.Printf("2. Run the installer\n")
	fmt.Printf("3. Restart your terminal\n")
	fmt.Printf("4. Run '%sintel start%s' again\n", output.GreenColor, output.Reset)
	
	return fmt.Errorf("manual installation required")
}

// startOllama attempts to start the Ollama service
func (om *OllamaManager) startOllama() error {
	fmt.Printf("%s🚀 Starting Ollama service...%s\n", output.YellowColor, output.Reset)
	
	// Try to start Ollama in the background
	cmd := exec.Command(om.binaryPath, "serve")
	
	// Start the process in the background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Ollama service: %w", err)
	}
	
	// Wait a moment for the service to start
	time.Sleep(2 * time.Second)
	
	// Check if the service is now running
	if err := om.checkService(); err != nil {
		return fmt.Errorf("Ollama service failed to start properly: %w", err)
	}
	
	fmt.Printf("%s✅ Ollama service started successfully%s\n", output.GreenColor, output.Reset)
	return nil
}

// GetStatus returns the current status of Ollama
func (om *OllamaManager) GetStatus() (string, error) {
	if err := om.checkInstallation(); err != nil {
		return "❌ Not installed", err
	}
	
	if err := om.checkService(); err != nil {
		return "⚠️  Installed but not running", err
	}
	
	return "✅ Running", nil
}

// GetBinaryPath returns the path to the Ollama binary
func (om *OllamaManager) GetBinaryPath() string {
	return om.binaryPath
}

// GetServiceURL returns the Ollama service URL
func (om *OllamaManager) GetServiceURL() string {
	return om.serviceURL
}

// IsAvailable returns true if Ollama is installed and running
func (om *OllamaManager) IsAvailable() bool {
	return om.isInstalled && om.isRunning
}

// downloadFile downloads a file from a URL to a local path
func (om *OllamaManager) downloadFile(url, filepath string) error {
	fmt.Printf("%s⬇️  Downloading %s...%s\n", output.CyanColor, url, output.Reset)
	
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	
	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Check status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	
	// Write to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	
	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filepath, 0755); err != nil {
			return err
		}
	}
	
	fmt.Printf("%s✅ Download completed%s\n", output.GreenColor, output.Reset)
	return nil
}

// ShowManualInstructions displays manual installation instructions
func (om *OllamaManager) ShowManualInstructions() {
	fmt.Printf("\n%s📖 Manual Installation Instructions:%s\n", output.BoldColor, output.Reset)
	fmt.Printf("1. Visit: %s%s%s\n", output.CyanColor, om.downloadURL, output.Reset)
	fmt.Printf("2. Download the appropriate installer for your OS\n")
	fmt.Printf("3. Install and restart your terminal\n")
	fmt.Printf("4. Run '%sintel start%s' again\n", output.GreenColor, output.Reset)
	fmt.Printf("\nAlternatively, you can install via package managers:\n")
	
	switch runtime.GOOS {
	case "linux":
		fmt.Printf("• Ubuntu/Debian: Use the official install script\n")
		fmt.Printf("• Arch Linux: %syay -S ollama%s\n", output.CyanColor, output.Reset)
	case "darwin":
		fmt.Printf("• Homebrew: %sbrew install ollama%s\n", output.CyanColor, output.Reset)
	case "windows":
		fmt.Printf("• Chocolatey: %schoco install ollama%s\n", output.CyanColor, output.Reset)
		fmt.Printf("• Scoop: %sscoop install ollama%s\n", output.CyanColor, output.Reset)
	}
}