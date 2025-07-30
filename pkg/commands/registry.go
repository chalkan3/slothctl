package commands

import (
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

// RegisterCommands registers all discovered BluePrintCommand implementations with the root command.
func RegisterCommands(rootCmd *cobra.Command) {
	// This map will hold all Cobra commands, keyed by their name.
	cobraCommands := make(map[string]*cobra.Command)

	// First pass: Create all Cobra commands and add them to the map.
	for _, cmd := range commands {
		cobraCmd := cmd.CobraCommand()
		cobraCommands[cobraCmd.Name()] = cobraCmd
	}

	// Second pass: Establish parent-child relationships.
	for _, cmd := range commands {
		cobraCmd := cobraCommands[cmd.CobraCommand().Name()]
		if cmd.Parent() != "" {
			parentCmd, ok := cobraCommands[cmd.Parent()]
			if ok {
				parentCmd.AddCommand(cobraCmd)
			} else {
				log.Warn("Parent command not found", "parent", cmd.Parent(), "command", cmd.CobraCommand().Name())
				rootCmd.AddCommand(cobraCmd)
			}
		} else {
			rootCmd.AddCommand(cobraCmd)
		}
	}
}

// AddCommandToRegistry is called by individual command packages' init() functions
// to register themselves with the command registry.
func AddCommandToRegistry(cmd BluePrintCommand) {
	commands = append(commands, cmd)
}

// LoadPlugins attempts to load commands from Go plugin files (.so).
// This is a more advanced feature and requires Go 1.8+ and building with `go build -buildmode=plugin`.
// For simplicity, this is a placeholder and not fully implemented for this example.
/*
func LoadPlugins(rootCmd *cobra.Command, pluginDir string) {
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		log.Error("Could not read plugin directory", "directory", pluginDir, "error", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			p, err := plugin.Open(pluginPath)
			if err != nil {
				log.Error("Failed to load plugin", "plugin_path", pluginPath, "error", err)
				continue
			}

			sym, err := p.Lookup("Command") // Assuming each plugin exports a 'Command' symbol
			if err != nil {
				log.Error("Failed to lookup 'Command' symbol in plugin", "plugin_path", pluginPath, "error", err)
				continue
			}

			cmd, ok := sym.(BluePrintCommand)
			if !ok {
				log.Error("Plugin does not implement BluePrintCommand interface", "plugin_path", pluginPath)
				continue
			}

			AddCommandToRegistry(cmd)
			log.Info("Loaded command from plugin", "command", cmd.CobraCommand().Name(), "plugin_path", pluginPath)
		}
	}
}
*/
