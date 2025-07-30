package vpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const DefaultConfigFile = "default"
const pidFileName = "openfortivpn.pid"
const logFileName = "openfortivpn.log"

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

// GetVPnPidFilePath returns the full path to the PID file.
func GetVPnPidFilePath() (string, error) {
	configDir, err := GetVPNConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, pidFileName), nil
}

// GetVPnLogFilePath returns the full path to the log file.
func GetVPnLogFilePath() (string, error) {
	configDir, err := GetVPNConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, logFileName), nil
}

// WriteVPnPid writes the PID to the PID file.
func WriteVPnPid(pid int) error {
	pidFilePath, err := GetVPnPidFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0644)
}

// ReadVPnPid reads the PID from the PID file.
func ReadVPnPid() (int, error) {
	pidFilePath, err := GetVPnPidFilePath()
	if err != nil {
		return 0, err
	}
	content, err := os.ReadFile(pidFilePath)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %w", err)
	}
	return pid, nil
}

// DeleteVPnPidFile removes the PID file.
func DeleteVPnPidFile() error {
	pidFilePath, err := GetVPnPidFilePath()
	if err != nil {
		return err
	}
	return os.Remove(pidFilePath)
}
