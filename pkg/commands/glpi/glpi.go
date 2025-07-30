package glpi

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// glpiCmd represents the base command for 'glpi'
type glpiCmd struct{}

func (c *glpiCmd) Parent() string {
	return ""
}

func (c *glpiCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "glpi",
		Short: "Manage GLPI instances and tickets",
		Long:  `The glpi command provides tools to manage GLPI instances, including registration, and interaction with tickets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&glpiCmd{})
}
