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

// registerCmd represents the 'server register' command
type registerCmd struct{}

func (c *registerCmd) Parent() string {
	return "server"
}

func (c *registerCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Short: "Registers a new server",
		Long:  `Registers a new server entry with its group, context, IP, and user details.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			group, _ := cmd.Flags().GetString("group")
			context, _ := cmd.Flags().GetString("context")
			ip, _ := cmd.Flags().GetString("ip")
			user, _ := cmd.Flags().GetString("user")

			if group == "" || context == "" || ip == "" || user == "" {
				return fmt.Errorf("group, context, ip, and user flags are required")
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

			server := servermanager.Server{
				Name:    name,
				Group:   group,
				Context: context,
				IP:      ip,
				User:    user,
			}

			if err := sm.SaveServer(server); err != nil {
				return fmt.Errorf("failed to save server: %w", err)
			}

			log.Info("Server registered successfully.", "name", name, "group", group, "context", context)
			return nil
		},
	}

	cmd.Flags().StringP("group", "g", "", "Server group (required)")
	cmd.Flags().StringP("context", "c", "", "Server context (required)")
	cmd.Flags().StringP("ip", "i", "", "Server IP address (required)")
	cmd.Flags().StringP("user", "u", "", "SSH username for the server (required)")

	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("context")
	cmd.MarkFlagRequired("ip")
	cmd.MarkFlagRequired("user")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&registerCmd{})
}
