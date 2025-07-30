package vpn

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
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

			pid, err := ReadVPnPid()
			if err != nil {
				log.Warn("Could not read VPN PID file. Attempting to check by process name.", "error", err)
				// Fallback to pgrep if PID file not found
				statusCmd := exec.Command("pgrep", "-a", "openfortivpn")
				output, err := statusCmd.CombinedOutput()

				if err != nil {
					log.Info("VPN Status: Disconnected (No openfortivpn process found)")
					return nil
				}
				log.Info("VPN Status: Connected (via pgrep)")
				log.Info("Process details:", "process", string(output))
			} else {
				// Check if the process with the PID from file is running
				_, err := os.FindProcess(pid) // Removed 'process' variable
				if err != nil {
					log.Warn("Process not found for PID from file. Attempting to check by process name.", "pid", pid, "error", err)
					// Fallback to pgrep if process not found
					statusCmd := exec.Command("pgrep", "-a", "openfortivpn")
					output, err := statusCmd.CombinedOutput()

					if err != nil {
						log.Info("VPN Status: Disconnected (No openfortivpn process found)")
						return nil
					}
					log.Info("VPN Status: Connected (via pgrep)")
					log.Info("Process details:", "process", string(output))
				} else { // Process found
					// This is a basic check, a more robust one would check cmdline
					statusCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
					output, err := statusCmd.CombinedOutput()
					if err == nil && strings.TrimSpace(string(output)) == "openfortivpn" {
						log.Info("VPN Status: Connected", "pid", pid)
					} else {
						log.Info("VPN Status: Disconnected (PID file exists but process is not openfortivpn or not running)", "pid", pid)
					}
				}
			}

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
