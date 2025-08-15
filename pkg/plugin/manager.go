package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"forgeai/pkg/sandbox"
)

// Manifest represents the plugin manifest file
type Manifest struct {
	Name      string   `json:"name"`
	Languages []string `json:"languages"`
}

// Executor is the interface that all language executors must implement
type Executor interface {
	// Execute runs the provided code in a sandboxed environment
	Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error)

	// ExecuteFile runs the provided file in a sandboxed environment
	ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error)

	// SupportedLanguages returns a list of supported languages
	SupportedLanguages() []string
}

// ExternalExecutor implements the Executor interface for external executables
type ExternalExecutor struct {
	binaryPath string
	languages  []string
}

// NewExternalExecutor creates a new ExternalExecutor
func NewExternalExecutor(binaryPath string, languages []string) *ExternalExecutor {
	return &ExternalExecutor{
		binaryPath: binaryPath,
		languages:  languages,
	}
}

// Execute runs the provided code using the external executable
func (e *ExternalExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Prepare the command
	cmd := exec.CommandContext(ctx, e.binaryPath, "execute", language, code)
	
	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return nil, fmt.Errorf("failed to execute code: %w", err)
	}
	
	// Parse the JSON output
	var result sandbox.ExecutionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}
	
	return &result, nil
}

// ExecuteFile runs the provided file using the external executable
func (e *ExternalExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Prepare the command
	cmd := exec.CommandContext(ctx, e.binaryPath, "execute-file", filePath)
	
	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return nil, fmt.Errorf("failed to execute file: %w", err)
	}
	
	// Parse the JSON output
	var result sandbox.ExecutionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}
	
	return &result, nil
}

// SupportedLanguages returns a list of supported languages
func (e *ExternalExecutor) SupportedLanguages() []string {
	return e.languages
}

// Manager handles plugin loading and management
type Manager struct {
	plugins map[string]Executor
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Executor),
	}
}

// LoadPlugin loads a plugin from the specified path
func (m *Manager) LoadPlugin(pluginDir string) error {
	// Read the manifest file
	manifestPath := filepath.Join(pluginDir, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	
	// Parse the manifest
	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	
	// Find the executable
	binaryPath := filepath.Join(pluginDir, manifest.Name)
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Try with .exe extension on Windows
		binaryPath = filepath.Join(pluginDir, manifest.Name+".exe")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			return fmt.Errorf("plugin executable not found: %s or %s.exe", manifest.Name, manifest.Name)
		}
	}
	
	// Create the executor
	executor := NewExternalExecutor(binaryPath, manifest.Languages)
	
	// Register the executor for each supported language
	for _, lang := range manifest.Languages {
		m.plugins[lang] = executor
	}
	
	return nil
}

// LoadPluginsFromDir loads all plugins from the specified directory
func (m *Manager) LoadPluginsFromDir(dir string) error {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory doesn't exist, that's okay
		return nil
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	// Load each plugin
	for _, entry := range entries {
		if entry.IsDir() {
			pluginDir := filepath.Join(dir, entry.Name())
			if err := m.LoadPlugin(pluginDir); err != nil {
				// Log error but continue loading other plugins
				fmt.Printf("Warning: failed to load plugin %s: %v\n", entry.Name(), err)
			}
		}
	}

	return nil
}

// GetExecutor returns the executor for the specified language
func (m *Manager) GetExecutor(language string) (Executor, bool) {
	executor, ok := m.plugins[language]
	return executor, ok
}

// SupportedLanguages returns a list of all supported languages from all plugins
func (m *Manager) SupportedLanguages() []string {
	langs := make([]string, 0, len(m.plugins))
	for lang := range m.plugins {
		langs = append(langs, lang)
	}
	return langs
}

// ListPlugins lists all plugins in the specified directory
func (m *Manager) ListPlugins(dir string) ([]string, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %w", err)
	}

	// Collect plugin names
	plugins := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}

	return plugins, nil
}