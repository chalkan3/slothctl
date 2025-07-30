package glpi

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// withCmd represents the 'glpi with' command
type withCmd struct{}

func (c *withCmd) Parent() string {
	return "glpi"
}

func (c *withCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "with [name]",
		Short: "Sets a GLPI instance as the default for subsequent commands",
		Long:  `Sets a registered GLPI instance as the default, so you don't need to specify it for every ticket command.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := glpimanager.NewManager(db)

			// Verify instance exists before setting as default
			_, err = sm.GetGLPIInstance(name)
			if err != nil {
				return fmt.Errorf("GLPI instance %s not found: %w", name, err)
			}

			if err := sm.SetDefaultGLPIInstance(name); err != nil {
				return fmt.Errorf("failed to set default GLPI instance: %w", err)
			}

			log.Info("Default GLPI instance set successfully.", "name", name)
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&withCmd{})
}
