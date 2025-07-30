package background

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// backgroundTaskCmd represents the base command for 'background-task'
type backgroundTaskCmd struct{}

func (c *backgroundTaskCmd) Parent() string {
	return ""
}

func (c *backgroundTaskCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "background-task",
		Short: "Manages background tasks",
		Long:  `The background-task command provides tools to start and manage long-running Go functions in the background.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&backgroundTaskCmd{})
}
