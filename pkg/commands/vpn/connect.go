package vpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// connectCmd represents the 'vpn connect' command
type connectCmd struct{}

func (c *connectCmd) Parent() string {
	return "vpn"
}

func (c *connectCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [config_name]",
		Short: "Connect to a VPN",
		Long:  `Starts a VPN connection using the specified configuration file or the default one if none is provided.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var configFile string
			passwordStdin, _ := cmd.Flags().GetBool("password-stdin")

			if len(args) == 0 {
				// No config file provided, try to use the default
				configDir, err := GetVPNConfigDir()
				if err != nil {
					return fmt.Errorf("could not get vpn config directory: %w", err)
				}
				defaultSymlinkPath := filepath.Join(configDir, DefaultConfigFile)

				linkTarget, err := os.Readlink(defaultSymlinkPath)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("no configuration file provided and no default VPN configuration set. Use 'slothctl vpn config set-default <name>' to set one.")
					} else {
						return fmt.Errorf("failed to read default VPN configuration symlink: %w", err)
					}
				}
				configFile = filepath.Join(configDir, linkTarget)
				log.Info("Using default VPN configuration.", "config", linkTarget)
			} else {
				// Config file provided as argument
				configName := args[0]
				configDir, err := GetVPNConfigDir()
				if err != nil {
					return fmt.Errorf("could not get vpn config directory: %w", err)
				}
				configFile = filepath.Join(configDir, configName)
			}

			log.Info("Attempting to connect to VPN...", "config_path", configFile)

			// Check if the config file exists
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", configFile)
			}

			// Construct the command to run openfortivpn
			var vpnCmd *exec.Cmd

			vpnArgs := []string{"-c", configFile}

			var password string
			if passwordStdin {
				log.Info("Reading password from stdin...")
				passwordBytes, err := os.ReadFile("/dev/stdin") // Read from stdin
				if err != nil {
					return fmt.Errorf("failed to read password from stdin: %w", err)
				}
				password = strings.TrimSpace(string(passwordBytes))
				vpnArgs = append(vpnArgs, "--password", password)
			}

			// Construct the full command string for bash -c
			// Using setsid to run in background and detach from terminal
			// Redirect stdin from /dev/null to prevent blocking
			// Redirect stdout/stderr to /dev/null to prevent output to terminal
			fullCmd := fmt.Sprintf("sudo setsid openfortivpn %s > /dev/null 2>&1 < /dev/null &", strings.Join(vpnArgs, " "))
			vpnCmd = exec.Command("bash", "-c", fullCmd)

			log.Info("Executing openfortivpn command in background using setsid...")
			if err := vpnCmd.Start(); err != nil { // Use Start() for background process
				log.Error("Failed to start VPN in background", "error", err, "output", err.Error()) // Use err.Error() for more details
				return fmt.Errorf("failed to start VPN: %w", err)
			}

			// Remove PID file if it exists, as we are not reliably tracking PID here
			// The status and disconnect commands will rely on pgrep/killall
			if err := DeleteVPnPidFile(); err != nil {
				log.Warn("Failed to delete old VPN PID file", "error", err)
			}

			log.Info("VPN connection process initiated in background.", "pid_of_bash_process", vpnCmd.Process.Pid) // Log PID of bash process for debugging
			log.Info("Use 'slothctl vpn status' to check the connection.")

			log.Info("VPN connection process started successfully.")
			log.Info("Use 'slothctl vpn status' to check the connection.")

			return nil
		},
	}
	cmd.Flags().Bool("password-stdin", false, "Read password from stdin for VPN authentication")
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&connectCmd{})
}
