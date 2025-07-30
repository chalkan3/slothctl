package server

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
)

// withCmd represents the 'server with' command
type withCmd struct{}

func (c *withCmd) Parent() string {
	return "server"
}

func (c *withCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "with [name]",
		Short: "Sets a server as the default for subsequent commands",
		Long:  `Sets a registered server as the default, so group and context flags are not needed for other commands.`,
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

			// Verify server exists before setting as default
			_, err = sm.GetServer(group, context, name)
			if err != nil {
				return fmt.Errorf("server %s:%s:%s not found: %w", group, context, name, err)
			}

			if err := sm.SetDefaultServer(group, context, name); err != nil {
				return fmt.Errorf("failed to set default server: %w", err)
			}

			log.Info("Default server set successfully.", "name", name, "group", group, "context", context)
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
	commands.AddCommandToRegistry(&withCmd{})
}
