package commands

import (
	"github.com/spf13/cobra"
	"github.com/chalkan3/slothctl/internal/log"
)

// BluePrintCommand defines the interface for a modular CLI command.
// Commands should implement this interface to be automatically registered.
type BluePrintCommand interface {
	// Parent returns the name of the parent command, or empty string if it's a root command.
	Parent() string
	// CobraCommand returns the Cobra command instance for this command.
	CobraCommand() *cobra.Command
}

// commands is a slice to hold all discovered BluePrintCommand implementations.
var commands []BluePrintCommand

// RegisterCommands registers all discovered BluePrintCommand implementations with the root command.
func RegisterCommands(rootCmd *cobra.Command) {
	cobraCommands := make(map[string]*cobra.Command)
	commandsToProcess := make(map[string]BluePrintCommand) // Store commands that need to be processed

	// First pass: Create all Cobra commands and add them to the map.
	// Also, populate commandsToProcess
	for _, cmd := range commands {
		cobraCmd := cmd.CobraCommand()
		cobraCommands[cobraCmd.Name()] = cobraCmd
		commandsToProcess[cobraCmd.Name()] = cmd
	}

	// Keep track of commands successfully added
	addedCommands := make(map[string]bool)

	// Repeatedly try to add commands until no more can be added in a pass
	for len(commandsToProcess) > 0 {
		commandsAddedInPass := 0
		for cmdName, cmd := range commandsToProcess {
			if cmd.Parent() == "" {
				// Root command
				rootCmd.AddCommand(cobraCommands[cmdName])
				addedCommands[cmdName] = true
				delete(commandsToProcess, cmdName)
				commandsAddedInPass++
			} else {
				parentCmd, ok := cobraCommands[cmd.Parent()]
				if ok && addedCommands[cmd.Parent()] { // Parent exists and has been added
					parentCmd.AddCommand(cobraCommands[cmdName])
					addedCommands[cmdName] = true
					delete(commandsToProcess, cmdName)
					commandsAddedInPass++
				}
			}
		}
		if commandsAddedInPass == 0 && len(commandsToProcess) > 0 {
			// No commands were added in this pass, but some remain.
			// This indicates a circular dependency or a missing parent.
			log.Warn("Could not register all commands due to missing parents or circular dependencies. Remaining commands:", "count", len(commandsToProcess))
			for cmdName, cmd := range commandsToProcess {
				log.Warn("  - Command:", "name", cmdName, "parent", cmd.Parent())
			}
			break // Exit loop to prevent infinite loop
		}
	}
}

// AddCommandToRegistry is called by individual command packages' init() functions
// to register themselves with the command registry.
func AddCommandToRegistry(cmd BluePrintCommand) {
	commands = append(commands, cmd)
}