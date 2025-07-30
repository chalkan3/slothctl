package vpn

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// configCmd represents the base command for 'vpn config'
type configCmd struct{}

func (c *configCmd) Parent() string {
	return "vpn"
}

func (c *configCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage VPN configuration files",
		Long:  `Provides subcommands to create, list, remove, and manage VPN configuration files.`,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&configCmd{})
}
