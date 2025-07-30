package server

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// getCmd represents the 'server get' command
type getCmd struct{}

func (c *getCmd) Parent() string {
	return "server"
}

func (c *getCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Retrieves details of a registered server",
		Long:  `Retrieves and displays the full details of a specific registered server.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			group, _ := cmd.Flags().GetString("group")
			context, _ := cmd.Flags().GetString("context")

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

			fmt.Printf("\n[%s]\n", server.Name)
			fmt.Printf("  ip -> %s\n", server.IP)
			fmt.Printf("  user -> %s\n", server.User)
			fmt.Printf("  group -> %s\n", server.Group)
			fmt.Printf("  context -> %s\n", server.Context)
			fmt.Printf("  description -> %s\n", server.Description)

			return nil
		},
	}

	cmd.Flags().StringP("group", "g", "", "Server group (optional, uses default if not provided)")
	cmd.Flags().StringP("context", "c", "", "Server context (optional, uses default if not provided)")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&getCmd{})
}
