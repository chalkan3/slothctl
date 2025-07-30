package vpn

import (
	"fmt"
	"os"
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

			pid, err := ReadVPnPid()
			if err != nil {
				log.Warn("Could not read VPN PID file. Attempting to kill by name.", "error", err)
				// Fallback to killall if PID file not found
				dcCmd := exec.Command("sudo", "killall", "openfortivpn")
				output, err := dcCmd.CombinedOutput()
				if err != nil {
					if len(output) > 0 && string(output) != "" {
						log.Warn("No running openfortivpn process found or an error occurred.", "output", string(output))
					} else {
						log.Info("No running openfortivpn process found.")
					}
					return nil
				}
			} else {
				log.Info("Killing VPN process by PID.", "pid", pid)
				process, err := os.FindProcess(pid)
				if err != nil {
					log.Warn("Process not found for PID from file. Attempting to kill by name.", "pid", pid, "error", err)
					// Fallback to killall if process not found
					dcCmd := exec.Command("sudo", "killall", "openfortivpn")
					output, err := dcCmd.CombinedOutput()
					if err != nil {
						if len(output) > 0 && string(output) != "" {
							log.Warn("No running openfortivpn process found or an error occurred.", "output", string(output))
						} else {
							log.Info("No running openfortivpn process found.")
						}
						return nil
					}
				} else {
					if err := process.Signal(os.Interrupt); err != nil { // Use os.Interrupt for graceful shutdown
						log.Warn("Failed to send interrupt signal, attempting to kill.", "error", err)
						if err := process.Kill(); err != nil {
							return fmt.Errorf("failed to kill VPN process: %w", err)
						}
					}
				}
			}

			// Clean up PID file
			if err := DeleteVPnPidFile(); err != nil {
				log.Warn("Failed to delete VPN PID file", "error", err)
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
