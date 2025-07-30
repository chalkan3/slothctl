package tickets

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpi"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// createCmd represents the 'glpi tickets create' command
type createCmd struct{}

func (c *createCmd) Parent() string {
	return "tickets"
}

func (c *createCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new GLPI ticket",
		Long:  `Creates a new GLPI ticket with specified details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			name, _ := cmd.Flags().GetString("name")
			content, _ := cmd.Flags().GetString("content")
			urgency, _ := cmd.Flags().GetInt("urgency")
			impact, _ := cmd.Flags().GetInt("impact")
			requesterID, _ := cmd.Flags().GetInt("requester-id")
			assigneeID, _ := cmd.Flags().GetInt("assignee-id")

			if name == "" || content == "" {
				return fmt.Errorf("name and content are required to create a ticket")
			}

			ticketInput := glpi.TicketInput{
				Name:        name,
				Content:     content,
				Urgency:     urgency,
				Impact:      impact,
				RequesterID: requesterID,
				AssigneeID:  assigneeID,
			}

			log.Info("Creating GLPI ticket...")
			createdTicket, err := client.CreateTicket(ticketInput)
			if err != nil {
				return fmt.Errorf("failed to create ticket: %w", err)
			}

			log.Info("Ticket created successfully.", "id", createdTicket.ID, "name", createdTicket.Name)
			fmt.Printf("Ticket ID: %d\n", createdTicket.ID)
			fmt.Printf("Name: %s\n", createdTicket.Name)
			fmt.Printf("Status: %s\n", glpi.GetStatusName(createdTicket.Status))
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name/title of the ticket (required)")
	cmd.Flags().StringP("content", "c", "", "Content/description of the ticket (required)")
	cmd.Flags().IntP("urgency", "u", 3, "Urgency of the ticket (1-5, 1=very low, 5=very high)")
	cmd.Flags().IntP("impact", "i", 3, "Impact of the ticket (1-5, 1=very low, 5=very high)")
	cmd.Flags().IntP("requester-id", "r", 0, "ID of the requester user")
	cmd.Flags().IntP("assignee-id", "a", 0, "ID of the assignee user")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("content")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&createCmd{})
}
