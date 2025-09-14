package saltnode

import (
	"os"
	"os/exec"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

type deleteCmd struct {
	cobraCommand *cobra.Command
}

func (c *deleteCmd) Parent() string {
	return "salt-node"
}

func (c *deleteCmd) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *deleteCmd) run(cmd *cobra.Command, args []string) {
	minionName := args[0]
	log.Info("Deleting salt node", "minion", minionName)

	cloneDir := "/tmp/salt-home"
	cloneURL := "https://github.com/chalkan3/salt-home"

	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		log.Info("Cloning repository", "url", cloneURL, "dir", cloneDir)
		gitCmd := exec.Command("git", "clone", cloneURL, cloneDir)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			log.Error("failed to clone repository", "error", err)
			return
		}
	} else {
		log.Info("Repository already cloned")
	}

	stackCmd := exec.Command("sudo", "PULUMI_CONFIG_PASSPHRASE=", "pulumi", "stack", "select", minionName)
	stackCmd.Dir = cloneDir
	stackCmd.Stdout = os.Stdout
	stackCmd.Stderr = os.Stderr
	if err := stackCmd.Run(); err != nil {
		log.Error("failed to select pulumi stack", "error", err)
		return
	}

	log.Info("Running pulumi destroy")
	pulumiCmd := exec.Command("sudo", "PULUMI_CONFIG_PASSPHRASE=", "pulumi", "destroy", "--yes", "--skip-preview")
	pulumiCmd.Dir = cloneDir
	pulumiCmd.Stdout = os.Stdout
	pulumiCmd.Stderr = os.Stderr

	if err := pulumiCmd.Run(); err != nil {
		log.Error("failed to run pulumi destroy", "error", err)
		return
	}

	log.Info("Successfully deleted salt node")
}

func NewDeleteCommand() commands.BluePrintCommand {
	deleteCmd := &deleteCmd{}

	cobraCmd := &cobra.Command{
		Use:   "delete <minion-name>",
		Short: "Deletes a salt node by running pulumi destroy.",
		Args:  cobra.ExactArgs(1),
		Run:   deleteCmd.run,
	}

	deleteCmd.cobraCommand = cobraCmd
	return deleteCmd
}

func init() {
	commands.AddCommandToRegistry(NewDeleteCommand())
}
