package vpn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// removeCmd represents the 'vpn config remove' command
type removeCmd struct{}

func (c *removeCmd) Parent() string {
	return "config"
}

func (c *removeCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [name]",
		Short: "Removes a VPN configuration file",
		Long:  `Removes a VPN configuration file from the slothctl config directory.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			configDir, err := GetVPNConfigDir()
			if err != nil {
				return fmt.Errorf("could not get vpn config directory: %w", err)
			}
			filePath := filepath.Join(configDir, name)

			log.Info("Attempting to remove VPN configuration...", "name", name, "path", filePath)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", filePath)
			}

			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove file: %w", err)
			}

			// Also remove default symlink if it points to this config
			defaultSymlinkPath := filepath.Join(configDir, DefaultConfigFile)
			linkTarget, err := os.Readlink(defaultSymlinkPath)
			if err == nil && linkTarget == name+".conf" {
				if err := os.Remove(defaultSymlinkPath); err != nil {
					log.Warn("Failed to remove default config symlink", "error", err)
				}
				log.Info("Removed default symlink as it pointed to the removed config.")
			}

			log.Info("Successfully removed VPN configuration.", "name", name)
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&removeCmd{})
}
