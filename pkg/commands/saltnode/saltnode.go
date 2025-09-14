package saltnode

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

type saltNodeCmd struct {
	cobraCommand *cobra.Command
}

func (c *saltNodeCmd) Parent() string {
	return ""
}

func (c *saltNodeCmd) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func NewSaltNodeCommand() commands.BluePrintCommand {
	cobraCmd := &cobra.Command{
		Use:   "salt-node",
		Short: "Manage salt nodes",
	}

	return &saltNodeCmd{
		cobraCommand: cobraCmd,
	}
}

func init() {
	commands.AddCommandToRegistry(NewSaltNodeCommand())
}
