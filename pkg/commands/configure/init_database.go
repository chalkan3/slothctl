package configure

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// databaseCmd represents the 'configure init database' command
type databaseCmd struct{}

func (c *databaseCmd) Parent() string {
	return "init"
}

func (c *databaseCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Initializes the embedded BoltDB database",
		Long:  `Initializes or ensures the embedded BoltDB database is set up correctly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Initializing embedded database...")

			// Initialize default config file (ensure config is loaded)
			if err := config.InitConfig(); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			// Reload config to ensure default values are loaded into AppConfig
			if err := config.LoadConfig(); err != nil {
				return fmt.Errorf("failed to reload config after initialization: %w", err)
			}
			log.Info("Database path from config", "path", config.AppConfig.DatabasePath)

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			log.Info("Expanded database path", "path", dbPath)

			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			// Create a default bucket if it doesn't exist
			err = db.Update(func(tx *bbolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte("slothctl_data"))
				if err != nil {
					return fmt.Errorf("create bucket: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to create default database bucket: %w", err)
			}

			log.Info("Embedded database initialized successfully.")
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&databaseCmd{})
}
