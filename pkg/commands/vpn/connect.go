package vpn

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
)

// connectCmd represents the 'vpn connect' command
type connectCmd struct{}

func (c *connectCmd) Parent() string {
	return "vpn"
}

func (c *connectCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [config_file]",
		Short: "Connect to a VPN",
		Long:  `Starts a VPN connection using the specified configuration file (e.g., an .ovpn file for OpenVPN).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := args[0]
			log.Info("Attempting to connect to VPN...", "config", configFile)

			// Check if the config file exists
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", configFile)
			}

			// Construct the command to run OpenVPN in the background
			// This requires running as root, so we use sudo.
			vpnCmd := exec.Command("sudo", "openvpn", "--config", configFile, "--daemon")

			// Capture and log output
			output, err := vpnCmd.CombinedOutput()
			if err != nil {
				log.Error("Failed to start VPN", "error", err, "output", string(output))
				return fmt.Errorf("failed to start VPN: %w", err)
			}

			log.Info("VPN connection process started successfully.", "output", string(output))
			log.Info("Use 'slothctl vpn status' to check the connection.")

			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&connectCmd{})
}
