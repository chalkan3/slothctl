package background

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// internalDaemonRunnerCmd represents the hidden command that runs the actual background task.
// It's not intended for direct user invocation.
type internalDaemonRunnerCmd struct{}

func (c *internalDaemonRunnerCmd) Parent() string {
	return ""
}

func (c *internalDaemonRunnerCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "internal-daemon-runner",
		Short:  "(Internal) Runs the background daemon task",
		Hidden: true, // Hide this command from help output
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Internal daemon runner started.", "pid", os.Getpid())

			// --- Daemon Logic Here ---
			// This is where your long-running Go function logic goes.
			// For demonstration, we'll just log a message every few seconds.

			// Set up signal handling for graceful shutdown
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

			// Loop indefinitely, or until a stop signal is received
			for i := 0; ; i++ {
				select {
				case sig := <-c:
					log.Info("Received signal, shutting down daemon gracefully.", "signal", sig.String())
					return nil // Exit the RunE function, which terminates the process
				case <-time.After(5 * time.Second):
					log.Info(fmt.Sprintf("Daemon heartbeat: %d seconds elapsed.", (i+1)*5))
				}
			}
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&internalDaemonRunnerCmd{})
}
