package server

import (
	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/pkg/commands"
)

// sshCmd represents the base command for 'server ssh'
type sshCmd struct{}

func (c *sshCmd) Parent() string {
	return "server"
}

func (c *sshCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Execute SSH operations on registered servers",
		Long:  `Provides subcommands to connect to or execute commands on registered servers via SSH.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&sshCmd{})
}
