package vpn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// createCmd represents the 'vpn config create' command
type createCmd struct{}

func (c *createCmd) Parent() string {
	return "config"
}

func (c *createCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Creates a new, managed VPN configuration file",
		Long:  `Generates a new VPN configuration file and saves it in the slothctl config directory.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			configDir, err := GetVPNConfigDir()
			if err != nil {
				return fmt.Errorf("could not get vpn config directory: %w", err)
			}
			filePath := filepath.Join(configDir, name)

			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetInt("port")
			user, _ := cmd.Flags().GetString("user")
			cert, _ := cmd.Flags().GetString("cert")

			log.Info("Creating new VPN configuration file...", "name", name)

			configContent := fmt.Sprintf(`host = %s
port = %d
username = %s
trusted-cert = %s
`, host, port, user, cert)

			if err := os.WriteFile(filePath, []byte(configContent), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			log.Info("Successfully created VPN configuration.", "name", name, "path", filePath)
			return nil
		},
	}

	cmd.Flags().String("host", "201.48.63.233", "VPN gateway host")
	cmd.Flags().Int("port", 10443, "VPN gateway port")
	cmd.Flags().String("user", "igor.rodrigues", "VPN username")
	cmd.Flags().String("cert", "9a820421292fb9a74a3b60f502532d2eb5bb539b7dd13046e912522d6ffde185", "Trusted certificate hash")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&createCmd{})
}
