package vpn

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
)

// statusCmd represents the 'vpn status' command
type statusCmd struct{}

func (c *statusCmd) Parent() string {
	return "vpn"
}

func (c *statusCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the VPN connection status",
		Long:  `Checks for a running OpenVPN process and displays the status of the tun0 network interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Checking VPN status...")

			// Check for a running openvpn process
			statusCmd := exec.Command("pgrep", "-a", "openvpn")
			output, err := statusCmd.CombinedOutput()

			if err != nil {
				log.Info("VPN Status: ${C_RED}Disconnected${C_RESET} (No OpenVPN process found)")
				return nil
			}

			log.Info("VPN Status: ${C_GREEN}Connected${C_RESET}")
			log.Info("Process details:", "process", string(output))

			// Check the status of the tun0 interface
			ipCmd := exec.Command("ip", "addr", "show", "tun0")
			ipOutput, err := ipCmd.CombinedOutput()
			if err != nil {
				log.Warn("Could not get status for tun0 interface. It might have a different name or not be up yet.")
			} else {
				fmt.Println("\n--- tun0 Interface Status ---")
				fmt.Println(string(ipOutput))
			}

			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&statusCmd{})
}
