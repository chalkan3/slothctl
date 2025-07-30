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
		Short: "Creates a sample openfortivpn configuration file",
		Long:  `Generates a new configuration file with a standard template for openfortivpn.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			log.Info("Creating new VPN configuration file...", "file", filename)

			configContent := `### configuration file for openfortivpn, see man openfortivpn(1) ###
#
trusted-cert = 9a820421292fb9a74a3b60f502532d2eb5bb539b7dd13046e912522d6ffde185
port = 10443
username = igor.rodrigues
host = 201.48.63.233
disallow-invalid-cert = 1
`

			if err := os.WriteFile(filename, []byte(configContent), 0644); err != nil {
				log.Error("Failed to write configuration file", "error", err)
				return fmt.Errorf("failed to write file: %w", err)
			}

			log.Info("Successfully created VPN configuration file.", "path", filename)
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&createConfigCmd{})
}
