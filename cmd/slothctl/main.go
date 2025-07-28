package main

import (
	"slothctl/internal/log"
	"slothctl/pkg/commands"
	"slothctl/pkg/config"
	_ "slothctl/zz_generated_commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slothctl",
	Short: "slothctl is a CLI tool for managing sloth-related tasks",
	Long:  `A powerful and flexible CLI tool to streamline your sloth-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.LoadConfig()
	},
}

func main() {
	commands.RegisterCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error executing slothctl", "error", err)
	}
}
