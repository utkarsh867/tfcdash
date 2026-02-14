// Package credentials provides functionality to read Terraform CLI credentials
// from the standard credentials file location (~/.terraform.d/credentials.tfrc.json)
package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CredentialsStore represents the structure of the Terraform CLI credentials file
type CredentialsStore struct {
	Credentials map[string]struct {
		Token string `json:"token"`
	} `json:"credentials"`
}

// DefaultCredentialsPath returns the default path to the Terraform CLI credentials file
func DefaultCredentialsPath() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// On Windows, use %APPDATA%
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			appdata = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		configDir = filepath.Join(appdata, "terraform.d")
	default:
		// On Unix-like systems (Linux, macOS), use ~/.terraform.d
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to HOME environment variable
			homeDir = os.Getenv("HOME")
		}
		configDir = filepath.Join(homeDir, ".terraform.d")
	}

	return filepath.Join(configDir, "credentials.tfrc.json")
}

// Load reads the Terraform CLI credentials from the specified path
func Load(path string) (*CredentialsStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var store CredentialsStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &store, nil
}

// LoadDefault reads the Terraform CLI credentials from the default location
func LoadDefault() (*CredentialsStore, error) {
	return Load(DefaultCredentialsPath())
}

// GetToken retrieves the API token for a specific hostname
func (c *CredentialsStore) GetToken(hostname string) (string, bool) {
	if c.Credentials == nil {
		return "", false
	}

	creds, exists := c.Credentials[hostname]
	if !exists {
		return "", false
	}

	return creds.Token, true
}

// GetDefaultToken retrieves the API token for Terraform Cloud (app.terraform.io)
func (c *CredentialsStore) GetDefaultToken() (string, bool) {
	// Try app.terraform.io first (Terraform Cloud)
	if token, ok := c.GetToken("app.terraform.io"); ok {
		return token, true
	}

	// Fall back to terraform.io
	if token, ok := c.GetToken("terraform.io"); ok {
		return token, true
	}

	return "", false
}

// ListHosts returns all hostnames that have stored credentials
func (c *CredentialsStore) ListHosts() []string {
	if c.Credentials == nil {
		return nil
	}

	hosts := make([]string, 0, len(c.Credentials))
	for host := range c.Credentials {
		hosts = append(hosts, host)
	}

	return hosts
}
