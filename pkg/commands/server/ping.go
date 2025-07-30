package server

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// pingCmd represents the 'server ping' command
type pingCmd struct{}

func (c *pingCmd) Parent() string {
	return "server"
}

func (c *pingCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping [name]",
		Short: "Pings a registered server",
		Long:  `Pings a registered server to check its reachability.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			group, _ := cmd.Flags().GetString("group")
			context, _ := cmd.Flags().GetString("context")

			if group == "" || context == "" {
				return fmt.Errorf("group and context flags are required")
			}

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := servermanager.NewManager(db)
			server, err := sm.GetServer(group, context, name)
			if err != nil {
				return fmt.Errorf("failed to get server: %w", err)
			}

			log.Info("Pinging server...", "name", server.Name, "ip", server.IP)

			pingCmd := exec.Command("ping", "-c", "4", server.IP)
			output, err := pingCmd.CombinedOutput()
			if err != nil {
				log.Error("Ping failed", "error", err, "output", string(output))
				return fmt.Errorf("ping to %s failed: %w", server.IP, err)
			}

			log.Info("Ping successful!", "output", string(output))
			return nil
		},
	}

	cmd.Flags().StringP("group", "g", "", "Server group (required)")
	cmd.Flags().StringP("context", "c", "", "Server context (required)")

	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("context")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&pingCmd{})
}
