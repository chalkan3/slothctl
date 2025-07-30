package vpn

import (
	"os/exec"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the 'vpn disconnect' command
type disconnectCmd struct{}

func (c *disconnectCmd) Parent() string {
	return "vpn"
}

func (c *disconnectCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from the VPN",
		Long:  `Terminates the running OpenVPN process, disconnecting from the VPN.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Attempting to disconnect from VPN...")

			// Find and kill the openvpn process. `killall` is a convenient way to do this.
			// This requires running as root.
			dcCmd := exec.Command("sudo", "killall", "openvpn")

			output, err := dcCmd.CombinedOutput()
			if err != nil {
				// killall returns a non-zero exit code if no process is found, which is not a fatal error for us.
				if len(output) > 0 && string(output) != "" {
					log.Warn("No running OpenVPN process found or an error occurred.", "output", string(output))
				} else {
					log.Info("No running OpenVPN process found.")
				}
				return nil // Don't treat "no process found" as an error
			}

			log.Info("VPN disconnected successfully.")
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&disconnectCmd{})
}
