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

// updateCommentCmd represents the 'glpi tickets update comment' command
type updateCommentCmd struct{}

func (c *updateCommentCmd) Parent() string {
	return "tickets"
}

func (c *updateCommentCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update comment [ticket-id]",
		Short: "Adds a comment/follow-up to a GLPI ticket",
		Long:  `Adds a comment or follow-up to a specified GLPI ticket. If no ticket ID is provided, the default ticket ID is used.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var ticketID int
			comment, _ := cmd.Flags().GetString("comment")

			if comment == "" {
				return fmt.Errorf("comment cannot be empty")
			}

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

			log.Info("Adding comment to ticket...", "ticket_id", ticketID)
			if err := client.AddTicketFollowup(ticketID, comment); err != nil {
				return fmt.Errorf("failed to add comment to ticket %d: %w", ticketID, err)
			}

			log.Info("Comment added successfully.", "ticket_id", ticketID)
			return nil
		},
	}

	cmd.Flags().StringP("comment", "m", "", "The comment/follow-up message (required)")
	cmd.MarkFlagRequired("comment")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&updateCommentCmd{})
}
