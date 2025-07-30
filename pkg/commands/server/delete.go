package server

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// deleteCmd represents the 'server delete' command
type deleteCmd struct{}

func (c *deleteCmd) Parent() string {
	return "server"
}

func (c *deleteCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Deletes a registered server",
		Long:  `Deletes a server entry based on its name, group, and context.`,
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
			if err := sm.Init(); err != nil {
				return fmt.Errorf("failed to initialize server manager: %w", err)
			}

			if err := sm.DeleteServer(group, context, name); err != nil {
				return fmt.Errorf("failed to delete server: %w", err)
			}

			log.Info("Server deleted successfully.", "name", name, "group", group, "context", context)
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
	commands.AddCommandToRegistry(&deleteCmd{})
}
