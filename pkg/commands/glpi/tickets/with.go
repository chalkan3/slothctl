package tickets

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// withCmd represents the 'glpi tickets with' command
type withCmd struct{}

func (c *withCmd) Parent() string {
	return "tickets"
}

func (c *withCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "with [id]",
		Short: "Sets a GLPI ticket ID as the default for subsequent commands",
		Long:  `Sets a GLPI ticket ID as the default, so you don't need to specify it for every ticket update command.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketIDStr := args[0]
			ticketID, err := strconv.Atoi(ticketIDStr)
			if err != nil {
				return fmt.Errorf("invalid ticket ID: %w", err)
			}

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := glpimanager.NewManager(db)

			// Optionally, verify ticket exists before setting as default
			// client, err := sm.GetDefaultGLPIClient()
			// if err != nil {
			// 	return fmt.Errorf("failed to get default GLPI client: %w", err)
			// }
			// _, err = client.GetTicket(ticketID)
			// if err != nil {
			// 	return fmt.Errorf("ticket %d not found: %w", ticketID, err)
			// }

			if err := sm.SetDefaultTicketID(ticketID); err != nil {
				return fmt.Errorf("failed to set default ticket ID: %w", err)
			}

			log.Info("Default GLPI ticket ID set successfully.", "id", ticketID)
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&withCmd{})
}
