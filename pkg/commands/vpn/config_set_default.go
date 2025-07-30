package vpn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// setDefaultCmd represents the 'vpn config set-default' command
type setDefaultCmd struct{}

func (c *setDefaultCmd) Parent() string {
	return "config"
}

func (c *setDefaultCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default [name]",
		Short: "Sets a VPN configuration as default",
		Long:  `Creates a symlink to the specified VPN configuration, making it the default.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			configDir, err := GetVPNConfigDir()
			if err != nil {
				return fmt.Errorf("could not get vpn config directory: %w", err)
			}
			filePath := filepath.Join(configDir, name+".conf")
			defaultSymlinkPath := filepath.Join(configDir, DefaultConfigFile)

			log.Info("Setting default VPN configuration...", "name", name)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", filePath)
			}

			// Remove existing default symlink if it exists
			if _, err := os.Lstat(defaultSymlinkPath); err == nil {
				if err := os.Remove(defaultSymlinkPath); err != nil {
					return fmt.Errorf("failed to remove existing default symlink: %w", err)
				}
			}

			// Create new symlink
			if err := os.Symlink(name+".conf", defaultSymlinkPath); err != nil {
				return fmt.Errorf("failed to create default symlink: %w", err)
			}

			log.Info("Successfully set default VPN configuration.", "name", name)
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&setDefaultCmd{})
}
