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

// updateOwnerCmd represents the 'glpi tickets update owner' command
type updateOwnerCmd struct{}

func (c *updateOwnerCmd) Parent() string {
	return "tickets"
}

func (c *updateOwnerCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update owner [ticket-id]",
		Short: "Updates the owner (assignee) of a GLPI ticket",
		Long:  `Updates the assigned owner of a specified GLPI ticket. If no ticket ID is provided, the default ticket ID is used.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ownerName, _ := cmd.Flags().GetString("owner")

			if ownerName == "" {
				return fmt.Errorf("owner name cannot be empty")
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

			// Get user ID from owner name
			users, err := client.GetUsers()
			if err != nil {
				return fmt.Errorf("failed to get users: %w", err)
			}
			ownerID := -1
			for _, user := range users {
				if strings.EqualFold(user.Name, ownerName) {
					ownerID = user.ID
					break
				}
			}
			if ownerID == -1 {
				return fmt.Errorf("owner user '%s' not found in GLPI", ownerName)
			}

			log.Info("Updating ticket owner...", "ticket_id", ticketID, "owner", ownerName)
			ticketInput := glpi.TicketInput{AssigneeID: ownerID}
			if err := client.UpdateTicket(ticketID, ticketInput); err != nil {
				return fmt.Errorf("failed to update ticket owner for ticket %d: %w", ticketID, err)
			}

			log.Info("Ticket owner updated successfully.", "ticket_id", ticketID, "new_owner", ownerName)
			return nil
		},
	}

	cmd.Flags().StringP("owner", "o", "", "The new owner (assignee) username for the ticket (required)")
	cmd.MarkFlagRequired("owner")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&updateOwnerCmd{})
}
