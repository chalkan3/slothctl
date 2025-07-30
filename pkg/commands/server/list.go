package server

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/servermanager"
)

// listCmd represents the 'server list' command
type listCmd struct{}

func (c *listCmd) Parent() string {
	return "server"
}

func (c *listCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all registered servers",
		Long:  `Lists all registered servers, grouped by context and group, with their details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			servers, err := sm.ListServers()
			if err != nil {
				return fmt.Errorf("failed to list servers: %w", err)
			}

			if len(servers) == 0 {
				log.Info("No servers registered.")
				return nil
			}

			// Group servers by group -> context
			groupedServers := make(map[string]map[string][]servermanager.Server)
			for _, s := range servers {
				if _, ok := groupedServers[s.Group]; !ok {
					groupedServers[s.Group] = make(map[string][]servermanager.Server)
				}
				groupedServers[s.Group][s.Context] = append(groupedServers[s.Group][s.Context], s)
			}

			// Sort groups and contexts for consistent output
			groups := make([]string, 0, len(groupedServers))
			for g := range groupedServers {
				groups = append(groups, g)
			}
			sort.Strings(groups)

			for _, g := range groups {
				fmt.Printf("[%s]\n", g)
				contexts := make([]string, 0, len(groupedServers[g]))
				for c := range groupedServers[g] {
					contexts = append(contexts, c)
				}
				sort.Strings(contexts)

				for _, c := range contexts {
					fmt.Printf("  [%s]\n", c)
					sort.Slice(groupedServers[g][c], func(i, j int) bool {
						return groupedServers[g][c][i].Name < groupedServers[g][c][j].Name
					})
					for _, s := range groupedServers[g][c] {
						fmt.Printf("    [%s]\n", s.Name)
						fmt.Printf("      ip -> %s\n", s.IP)
						fmt.Printf("      user -> %s\n", s.User)
						fmt.Printf("      description -> %s\n", s.Description)
					}
				}
			}

			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&listCmd{})
}
