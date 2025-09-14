//go:build linux

package background

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

// startCmd represents the 'background-task start' command
type startCmd struct{}

func (c *startCmd) Parent() string {
	return "background-task"
}

func (c *startCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a Go function in the background as a detached process",
		Long:  `This command launches a specified Go function as a new, detached process that continues to run after the main CLI command exits.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Attempting to start background task...")

			executable, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}

			// Define the command to run the internal daemon runner
			// We pass a special argument to indicate it's the internal runner
			commandArgs := []string{"internal-daemon-runner"}

			// Prepare the command
			proc := exec.Command(executable, commandArgs...)

			// Detach the process from the current session
			proc.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

			// Redirect standard I/O to /dev/null or log files
			// For a real daemon, you'd want to redirect to log files.
			// For simplicity, redirecting to /dev/null for now.
			proc.Stdin = nil // No stdin for daemon

			// Open log files for stdout and stderr
			logDir := filepath.Join(os.TempDir(), "slothctl_daemon_logs")
			if err := os.MkdirAll(logDir, 0755); err != nil {
				return fmt.Errorf("failed to create log directory: %w", err)
			}
			stdoutLogPath := filepath.Join(logDir, "daemon_stdout.log")
			stderrLogPath := filepath.Join(logDir, "daemon_stderr.log")

			stdoutFile, err := os.OpenFile(stdoutLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open stdout log file: %w", err)
			}
			defer stdoutFile.Close()

			stderrFile, err := os.OpenFile(stderrLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open stderr log file: %w", err)
			}
			defer stderrFile.Close()

			proc.Stdout = stdoutFile
			proc.Stderr = stderrFile

			// Start the detached process
			err = proc.Start()
			if err != nil {
				return fmt.Errorf("failed to start detached process: %w", err)
			}

			log.Info("Background task started successfully.", "pid", proc.Process.Pid)
			log.Info(fmt.Sprintf("Logs can be found in %s", logDir))
			return nil
		},
	}
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&startCmd{})
}
