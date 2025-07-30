package vpn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// listCmd represents the 'vpn config list' command
type listCmd struct{}

func (c *listCmd) Parent() string {
	return "config"
}

func (c *listCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists available VPN configuration files",
		Long:  `Lists all VPN configuration files managed by slothctl.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir, err := GetVPNConfigDir()
			if err != nil {
				return fmt.Errorf("could not get vpn config directory: %w", err)
			}

			log.Info("Listing VPN configurations from:", "path", configDir)

			files, err := os.ReadDir(configDir)
			if err != nil {
				return fmt.Errorf("failed to read config directory: %w", err)
			}

			if len(files) == 0 {
				log.Info("No VPN configurations found.")
				return nil
			}

			fmt.Println("\nAvailable VPN Configurations:")
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".conf" {
					fmt.Printf("- %s\n", file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]) // Remove .conf extension
				}
			}

			// Check for default config
			defaultConfigPath := filepath.Join(configDir, DefaultConfigFile)
			defaultTarget, err := os.Readlink(defaultConfigPath)
			if err == nil {
				fmt.Printf("\nDefault VPN Configuration: %s\n", filepath.Base(defaultTarget))
			} else if !os.IsNotExist(err) {
				log.Warn("Could not read default config symlink", "error", err)
			}

			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&listCmd{})
}
