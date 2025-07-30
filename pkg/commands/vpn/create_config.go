package vpn

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
)

// createConfigCmd represents the 'vpn create-config' command
type createConfigCmd struct{}

func (c *createConfigCmd) Parent() string {
	return "vpn"
}

func (c *createConfigCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-config [filename]",
		Short: "Creates a customizable openfortivpn configuration file",
		Long:  `Generates a new configuration file for openfortivpn, allowing customization via flags.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetInt("port")
			user, _ := cmd.Flags().GetString("user")
			cert, _ := cmd.Flags().GetString("cert")

			log.Info("Creating new VPN configuration file...", "file", filename)

			configContent := fmt.Sprintf(`### configuration file for openfortivpn, see man openfortivpn(1) ###
#
host = %s
port = %d
username = %s
trusted-cert = %s
disallow-invalid-cert = 1
`, host, port, user, cert)

			if err := os.WriteFile(filename, []byte(configContent), 0644); err != nil {
				log.Error("Failed to write configuration file", "error", err)
				return fmt.Errorf("failed to write file: %w", err)
			}

			log.Info("Successfully created VPN configuration file.", "path", filename)
			return nil
		},
	}

	// Add flags for customization
	cmd.Flags().String("host", "201.48.63.233", "VPN gateway host")
	cmd.Flags().Int("port", 10443, "VPN gateway port")
	cmd.Flags().String("user", "igor.rodrigues", "VPN username")
	cmd.Flags().String("cert", "9a820421292fb9a74a3b60f502532d2eb5bb539b7dd13046e912522d6ffde185", "Trusted certificate hash")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&createConfigCmd{})
}