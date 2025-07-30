package tickets

import (
	"fmt"
	"os"
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

// listCmd represents the 'glpi tickets list' command
type listCmd struct{}

func (c *listCmd) Parent() string {
	return "tickets"
}

func (c *listCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists GLPI tickets",
		Long:  `Lists GLPI tickets, with options to filter by status or owner.`,
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

			statusStr, _ := cmd.Flags().GetString("status")
			ownerName, _ := cmd.Flags().GetString("by-owner")

			var statusIDs []int
			if statusStr != "" {
				statuses := strings.Split(statusStr, ",")
				for _, s := range statuses {
					id, err := glpi.GetStatusID(strings.ToLower(strings.TrimSpace(s)))
					if err != nil {
						return err
					}
					statusIDs = append(statusIDs, id)
				}
			}

			tickets, err := client.ListTickets(statusIDs)
			if err != nil {
				return fmt.Errorf("failed to list tickets: %w", err)
			}

			if ownerName != "" {
				users, err := client.GetUsers()
				if err != nil {
					return fmt.Errorf("failed to get users for owner filtering: %w", err)
				}
				ownerID := -1
				for _, user := range users {
					if strings.EqualFold(user.Name, ownerName) {
						ownerID = user.ID
						break
					}
				}
				if ownerID == -1 {
					return fmt.Errorf("owner %s not found", ownerName)
				}

				filteredTickets := []glpi.Ticket{}
				for _, ticket := range tickets {
					if ticket.AssigneeID == ownerID || ticket.RequesterID == ownerID {
						filteredTickets = append(filteredTickets, ticket)
					}
				}
				tickets = filteredTickets
			}

			if len(tickets) == 0 {
				log.Info("No tickets found.")
				return nil
			}

			fmt.Println("\nGLPI Tickets:")
			for _, ticket := range tickets {
				fmt.Printf("  ID: %d\n", ticket.ID)
				fmt.Printf("  Name: %s\n", ticket.Name)
				fmt.Printf("  Status: %s\n", glpi.GetStatusName(ticket.Status))
				fmt.Printf("  Description: %s\n", ticket.Content)
				fmt.Printf("  --------------------\n")
			}

			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "", "Comma-separated list of statuses (e.g., new,assigned,solved)")
	cmd.Flags().StringP("by-owner", "o", "", "Filter tickets by owner name")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&listCmd{})
}
