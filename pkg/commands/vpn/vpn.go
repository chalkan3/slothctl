package vpn

import (
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// vpnCmd represents the base command for 'vpn'
type vpnCmd struct{}

func (c *vpnCmd) Parent() string {
	return ""
}

func (c *vpnCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpn",
		Short: "Manage VPN connections",
		Long:  `The vpn command provides tools to manage VPN connections, such as OpenVPN or WireGuard.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// By default, show help if no subcommand is given
			return cmd.Help()
		},
		TraverseChildren: true,
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&vpnCmd{})
}
