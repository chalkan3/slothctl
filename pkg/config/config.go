package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"slothctl/internal/log"
)

// Config represents the application configuration.
type Config struct {
	// AsdfInstallPath is the installation path for asdf.
	AsdfInstallPath string `mapstructure:"asdf_install_path"`
	// DatabasePath is the path to the embedded database file.
	DatabasePath    string `mapstructure:"database_path"`
}

// Global configuration instance.
var AppConfig Config

// LoadConfig initializes and loads the application configuration.
func LoadConfig() error {
	// Set default values for AppConfig directly
	AppConfig.AsdfInstallPath = os.ExpandEnv("$HOME/.asdf")
	AppConfig.DatabasePath = filepath.Join(os.ExpandEnv("$HOME/.slothctl"), "slothctl.db")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // type of config file
	viper.AddConfigPath("$HOME/.slothctl") // path to look for the config file in
	viper.AddConfigPath(".")             // optionally look for config in the working directory

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error, defaults are already set
			log.Info("Config file not found, using defaults or environment variables.")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal the config into the AppConfig struct, overriding defaults if present in config file
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}

	return nil
}

// InitConfig creates a default config file if it doesn't exist.
func InitConfig() error {
	configDir := os.ExpandEnv("$HOME/.slothctl")
	log.Info("Config directory", "path", configDir)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Info("Creating config directory", "path", configDir)
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	configPath := filepath.Join(configDir, "config.yaml")
	log.Info("Config file path", "path", configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Info("Creating default config file", "path", configPath)
		// Write current AppConfig values to file
		err := viper.WriteConfigAs(configPath)
		if err != nil {
			return fmt.Errorf("failed to write default config file: %w", err)
		}
		fmt.Printf("Default config file created at %s\n", configPath)
	}

	return nil
}
