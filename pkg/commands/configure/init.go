package configure

import (
	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/pkg/commands"
)

// initCmd represents the 'configure init' command
type initCmd struct{}

func (c *initCmd) Parent() string {
	return "configure"
}

func (c *initCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes slothctl components",
		Long:  "The init command provides subcommands to initialize various slothctl components, such as the database or control plane.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help() // Show help by default
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&initCmd{})
}