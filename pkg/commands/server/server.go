package server

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// serverCmd represents the base command for 'server'
type serverCmd struct{}

func (c *serverCmd) Parent() string {
	return ""
}

func (c *serverCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage server connections and configurations",
		Long:  `The server command provides tools to manage server entries, including registration, deletion, and SSH operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&serverCmd{})
}
