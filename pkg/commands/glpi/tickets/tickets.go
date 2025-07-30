package tickets

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// ticketsCmd represents the base command for 'glpi tickets'
type ticketsCmd struct{}

func (c *ticketsCmd) Parent() string {
	return "glpi"
}

func (c *ticketsCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tickets",
		Short: "Manage GLPI tickets",
		Long:  `The tickets command provides tools to interact with GLPI tickets, including listing, creating, updating, and retrieving details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&ticketsCmd{})
}
