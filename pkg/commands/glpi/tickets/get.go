package tickets

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpi"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// getCmd represents the 'glpi tickets get' command
type getCmd struct{}

func (c *getCmd) Parent() string {
	return "tickets"
}

func (c *getCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [ticket-id]",
		Short: "Retrieves details of a GLPI ticket",
		Long:  `Retrieves and displays the full details of a specific GLPI ticket. If no ticket ID is provided, the default ticket ID is used.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			log.Info("Retrieving ticket details...", "ticket_id", ticketID)
			ticket, err := client.GetTicket(ticketID)
			if err != nil {
				return fmt.Errorf("failed to retrieve ticket %d: %w", ticketID, err)
			}

			fmt.Printf("\nTicket Details (ID: %d):\n", ticket.ID)
			fmt.Printf("  Name: %s\n", ticket.Name)
			fmt.Printf("  Content: %s\n", ticket.Content)
			fmt.Printf("  Status: %s (%d)\n", glpi.GetStatusName(ticket.Status), ticket.Status)
			fmt.Printf("  Urgency: %d\n", ticket.Urgency)
			fmt.Printf("  Impact: %d\n", ticket.Impact)
			fmt.Printf("  Requester ID: %d\n", ticket.RequesterID)
			fmt.Printf("  Assignee ID: %d\n", ticket.AssigneeID)
			// Add more fields as needed

			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&getCmd{})
}
