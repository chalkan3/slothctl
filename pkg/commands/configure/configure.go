package configure

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// configureCmd represents the base command for 'configure'
type configureCmd struct{}

func (c *configureCmd) Parent() string {
	return ""
}

func (c *configureCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Manage slothctl configuration",
		Long:  `The configure command helps in setting up and managing various configurations for slothctl.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil // This command doesn't have a direct action, it's a parent for subcommands.
		},
		TraverseChildren: true,
	}
	// No flags for the base configure command
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&configureCmd{})
}
