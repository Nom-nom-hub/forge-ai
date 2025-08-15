package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"forgeai/pkg/plugin"
)

// PluginInfo represents metadata about a plugin
type PluginInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
	Languages   []string `json:"languages"`
	DownloadURL string   `json:"download_url"`
	FileHash    string   `json:"file_hash"`
	Signature   string   `json:"signature"`
}

// RegistryClient manages communication with the plugin registry
type RegistryClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewRegistryClient creates a new registry client
func NewRegistryClient(baseURL string) *RegistryClient {
	return &RegistryClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListPlugins retrieves a list of available plugins
func (rc *RegistryClient) ListPlugins() ([]PluginInfo, error) {
	url := fmt.Sprintf("%s/v1/plugins", rc.BaseURL)
	
	resp, err := rc.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plugins: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var plugins []PluginInfo
	if err := json.Unmarshal(body, &plugins); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return plugins, nil
}

// GetPlugin retrieves information about a specific plugin
func (rc *RegistryClient) GetPlugin(name string) (*PluginInfo, error) {
	url := fmt.Sprintf("%s/v1/plugins/%s", rc.BaseURL, name)
	
	resp, err := rc.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plugin: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var plugin PluginInfo
	if err := json.Unmarshal(body, &plugin); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &plugin, nil
}

// DownloadPlugin downloads a plugin to the specified directory
func (rc *RegistryClient) DownloadPlugin(name, version, destDir string) error {
	// Get plugin information
	pluginInfo, err := rc.GetPlugin(name)
	if err != nil {
		return fmt.Errorf("failed to get plugin info: %w", err)
	}
	
	// Create destination directory
	pluginDir := filepath.Join(destDir, name)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}
	
	// Download the plugin binary
	binaryURL := pluginInfo.DownloadURL
	if binaryURL == "" {
		binaryURL = fmt.Sprintf("%s/v1/plugins/%s/versions/%s/download", rc.BaseURL, name, version)
	}
	
	resp, err := rc.HTTPClient.Get(binaryURL)
	if err != nil {
		return fmt.Errorf("failed to download plugin binary: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	
	// Save the binary
	binaryName := pluginInfo.Name
	if filepath.Ext(binaryName) == "" {
		// Add appropriate extension based on OS
		if os.PathSeparator == '\\' {
			binaryName += ".exe"
		}
	}
	
	binaryPath := filepath.Join(pluginDir, binaryName)
	binaryFile, err := os.Create(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create binary file: %w", err)
	}
	defer binaryFile.Close()
	
	// Copy the binary data
	if _, err := io.Copy(binaryFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save binary: %w", err)
	}
	
	// Set executable permissions
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}
	
	// Create the manifest file
	manifest := plugin.Manifest{
		Name:      pluginInfo.Name,
		Languages: pluginInfo.Languages,
	}
	
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}
	
	manifestPath := filepath.Join(pluginDir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	
	return nil
}

// PluginManager manages local plugin installation and registry interaction
type PluginManager struct {
	LocalDir      string
	Registry      *RegistryClient
	PluginManager *plugin.Manager
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(localDir, registryURL string) *PluginManager {
	return &PluginManager{
		LocalDir:      localDir,
		Registry:      NewRegistryClient(registryURL),
		PluginManager: plugin.NewManager(),
	}
}

// InstallPlugin installs a plugin from the registry
func (pm *PluginManager) InstallPlugin(name, version string) error {
	return pm.Registry.DownloadPlugin(name, version, pm.LocalDir)
}

// ListInstalledPlugins lists locally installed plugins
func (pm *PluginManager) ListInstalledPlugins() ([]string, error) {
	return pm.PluginManager.ListPlugins(pm.LocalDir)
}

// ListRegistryPlugins lists available plugins from the registry
func (pm *PluginManager) ListRegistryPlugins() ([]PluginInfo, error) {
	return pm.Registry.ListPlugins()
}

// UpdatePlugin updates an installed plugin
func (pm *PluginManager) UpdatePlugin(name string) error {
	// For simplicity, we'll just reinstall the plugin
	// In a real implementation, we would check versions and only update if needed
	return pm.InstallPlugin(name, "latest")
}

// RemovePlugin removes an installed plugin
func (pm *PluginManager) RemovePlugin(name string) error {
	pluginDir := filepath.Join(pm.LocalDir, name)
	return os.RemoveAll(pluginDir)
}