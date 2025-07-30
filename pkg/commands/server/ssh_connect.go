package server

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
)

// sshConnectCmd represents the 'server ssh connect' command
type sshConnectCmd struct{}

func (c *sshConnectCmd) Parent() string {
	return "ssh"
}

func (c *sshConnectCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [name]",
		Short: "Connects to a registered server via SSH",
		Long:  `Establishes an SSH connection to a registered server. Can use password from stdin.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			group, _ := cmd.Flags().GetString("group")
			context, _ := cmd.Flags().GetString("context")
			passwordStdin, _ := cmd.Flags().GetBool("password-stdin")

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := servermanager.NewManager(db)

			// Resolve server details (using default if group/context not provided)
			if group == "" || context == "" {
				defaultGroup, defaultContext, defaultName, err := sm.GetDefaultServer()
				if err == nil && defaultName == name {
					group = defaultGroup
					context = defaultContext
				} else {
					return fmt.Errorf("group and context flags are required unless server is set as default")
				}
			}

			server, err := sm.GetServer(group, context, name)
			if err != nil {
				return fmt.Errorf("failed to get server: %w", err)
			}

			log.Info("Connecting to server via SSH...", "user", server.User, "ip", server.IP)

			sshArgs := []string{fmt.Sprintf("%s@%s", server.User, server.IP)}
			if passwordStdin {
				sshArgs = append(sshArgs, "-o", "PreferredAuthentications=password", "-o", "PubkeyAuthentication=no")
			}

			sshCmd := exec.Command("ssh", sshArgs...)

			// Connect stdin/stdout/stderr directly for interactive SSH
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr

			if passwordStdin {
				// Read password from stdin and pass to sshpass
				// This requires sshpass to be installed on the system.
				log.Info("Reading password from stdin...")
				passwordBytes, err := os.ReadFile("/dev/stdin") // Read from stdin
				if err != nil {
					return fmt.Errorf("failed to read password from stdin: %w", err)
				}
				password := strings.TrimSpace(string(passwordBytes))

				sshpassCmd := exec.Command("sshpass", "-p", password, "ssh", sshArgs...)
				sshpassCmd.Stdin = os.Stdin
				sshpassCmd.Stdout = os.Stdout
				sshpassCmd.Stderr = os.Stderr

				log.Info("Executing sshpass command...")
				if err := sshpassCmd.Run(); err != nil {
					return fmt.Errorf("sshpass command failed: %w", err)
				}
			} else {
				log.Info("Executing ssh command...")
				if err := sshCmd.Run(); err != nil {
					return fmt.Errorf("ssh command failed: %w", err)
				}
			}

			log.Info("SSH connection closed.")
			return nil
		},
	}

	cmd.Flags().StringP("group", "g", "", "Server group (optional, uses default if not provided)")
	cmd.Flags().StringP("context", "c", "", "Server context (optional, uses default if not provided)")
	cmd.Flags().Bool("password-stdin", false, "Read password from stdin for SSH authentication (requires sshpass)")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&sshConnectCmd{})
}
