package vpn

import (
	"os"
	"path/filepath"
)

const DefaultConfigFile = "default"

// GetVPNConfigDir returns the directory where VPN configurations are stored.
// It creates the directory if it doesn't exist.
func GetVPNConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "slothctl", "vpn")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}
