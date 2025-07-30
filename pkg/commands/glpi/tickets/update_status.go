package tickets

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpi"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// updateStatusCmd represents the 'glpi tickets update status' command
type updateStatusCmd struct{}

func (c *updateStatusCmd) Parent() string {
	return "tickets"
}

func (c *updateStatusCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update status [ticket-id]",
		Short: "Updates the status of a GLPI ticket",
		Long:  `Updates the status of a specified GLPI ticket. If no ticket ID is provided, the default ticket ID is used. Use --list-status to see available statuses.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			listStatus, _ := cmd.Flags().GetBool("list-status")
			statusStr, _ := cmd.Flags().GetString("status")

			if listStatus {
				fmt.Println("Available statuses: new, assigned, planned, pending, solved, closed")
				return nil
			}

			if statusStr == "" {
				return fmt.Errorf("status cannot be empty. Use --list-status to see available options.")
			}

			statusID, err := glpi.GetStatusID(strings.ToLower(strings.TrimSpace(statusStr)))
			if err != nil {
				return err
			}

			var ticketID int
			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := glpimanager.NewManager(db)
			client, err := sm.GetDefaultGLPIClient()
			if err != nil {
				return fmt.Errorf("failed to get default GLPI client: %w", err)
			}

			if len(args) > 0 {
				ticketID, err = strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("invalid ticket ID: %w", err)
				}
			} else {
				ticketID, err = sm.GetDefaultTicketID()
				if err != nil {
					return fmt.Errorf("no ticket ID provided and no default ticket ID set: %w", err)
				}
			}

			log.Info("Updating ticket status...", "ticket_id", ticketID, "status", statusStr)
			ticketInput := glpi.TicketInput{Status: statusID}
			if err := client.UpdateTicket(ticketID, ticketInput); err != nil {
				return fmt.Errorf("failed to update ticket status for ticket %d: %w", ticketID, err)
			}

			log.Info("Ticket status updated successfully.", "ticket_id", ticketID, "new_status", statusStr)
			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "", "The new status for the ticket (e.g., new, assigned, solved)")
	cmd.Flags().Bool("list-status", false, "List available statuses")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&updateStatusCmd{})
}
