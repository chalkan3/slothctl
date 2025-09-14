//go:build !linux

package background

import (
	"fmt"

	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// startCmd represents the 'background-task start' command
type startCmd struct{}

func (c *startCmd) Parent() string {
	return "background-task"
}

func (c *startCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a Go function in the background as a detached process (Linux only)",
		Long:  `This command is only supported on Linux.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("background tasks are only supported on Linux")
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&startCmd{})
}
