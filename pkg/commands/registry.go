package commands

import (
	"fmt"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/spf13/cobra"
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

// getUniqueCommandKey recursively constructs a unique key for a command based on its parentage.
// This is crucial to avoid name collisions for commands with the same name but different parents.
func getUniqueCommandKey(cmd BluePrintCommand, allBpCommands map[string]BluePrintCommand) string {
	if cmd.Parent() == "" {
		return cmd.CobraCommand().Name()
	}

	parentBpCmd, ok := allBpCommands[cmd.Parent()]
	if !ok {
		// This should ideally not happen if all commands are properly registered.
		// Fallback to just the command name if parent is not found, but log a warning.
		log.Warn("Parent not found for unique key generation", "parent", cmd.Parent(), "child", cmd.CobraCommand().Name())
		return cmd.CobraCommand().Name()
	}

	return fmt.Sprintf("%s/%s", getUniqueCommandKey(parentBpCmd, allBpCommands), cmd.CobraCommand().Name())
}

// RegisterCommands registers all discovered BluePrintCommand implementations with the root command.
func RegisterCommands(rootCmd *cobra.Command) {
	// Map to hold all BluePrintCommand instances, keyed by their simple name.
	// This is used to resolve parent BluePrintCommand objects.
	bpCommandsMap := make(map[string]BluePrintCommand)
	for _, cmd := range commands {
		bpCommandsMap[cmd.CobraCommand().Name()] = cmd
	}

	// This map will store all Cobra command instances, keyed by their unique full path.
	// Each Cobra command is created only once here.
	allCobraCommands := make(map[string]*cobra.Command)

	// First Loop: Initialize all Cobra command instances and store them in the map using unique keys.
	for _, cmd := range commands {
		uniqueKey := getUniqueCommandKey(cmd, bpCommandsMap)
		allCobraCommands[uniqueKey] = cmd.CobraCommand()
		log.Debug("Created Cobra command instance", "name", cmd.CobraCommand().Name(), "unique_key", uniqueKey)
	}

	// Second Loop: Establish parent-child relationships and add to rootCmd.
	for _, cmd := range commands {
		uniqueKey := getUniqueCommandKey(cmd, bpCommandsMap)
		currentCobraCmd := allCobraCommands[uniqueKey] // Retrieve the single instance

		if cmd.Parent() == "" {
			// It's a root-level command, add it directly to rootCmd.
			rootCmd.AddCommand(currentCobraCmd)
			log.Debug("Added root command", "name", currentCobraCmd.Name())
		} else {
			// It's a subcommand, find its parent and add it.
			parentUniqueKey := getUniqueCommandKey(bpCommandsMap[cmd.Parent()], bpCommandsMap)
			parentCobraCmd, ok := allCobraCommands[parentUniqueKey]
			if ok {
				parentCobraCmd.AddCommand(currentCobraCmd)
				log.Debug("Successfully added subcommand", "subcommand", currentCobraCmd.Name(), "parent", parentCobraCmd.Name())
			} else {
				// This case indicates a problem: a child command's parent was not found.
				// This could happen if the parent command itself wasn't registered, or if there's a typo.
				log.Warn("Parent Cobra command not found in map, adding child to root as fallback", "parent_unique_key", parentUniqueKey, "child", currentCobraCmd.Name())
				rootCmd.AddCommand(currentCobraCmd) // Fallback: add to root
			}
		}
	}
}

// AddCommandToRegistry is called by individual command packages' init() functions
// to register themselves with the command registry.
func AddCommandToRegistry(cmd BluePrintCommand) {
	commands = append(commands, cmd)
	log.Debug("Registering command (BluePrintCommand)", "name", cmd.CobraCommand().Name(), "parent", cmd.Parent())
}
