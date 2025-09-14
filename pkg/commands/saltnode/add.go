package saltnode

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/spf13/cobra"
)

type addCmd struct {
	cobraCommand *cobra.Command
	masterHost   string
	minionTarget string
	grains       []string
}

func (c *addCmd) Parent() string {
	return "salt-node"
}

func (c *addCmd) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *addCmd) run(cmd *cobra.Command, args []string) {
	minionName := args[0]
	log.Info("Adding new salt node", "minion", minionName, "master", c.masterHost, "target", c.minionTarget)

	cloneDir := "/tmp/salt-home"
	cloneURL := "https://github.com/chalkan3/salt-home"

	log.Info("Cloning repository", "url", cloneURL, "dir", cloneDir)
	if err := os.RemoveAll(cloneDir); err != nil {
		log.Error("failed to remove existing clone directory", "error", err)
		return
	}

	gitCmd := exec.Command("git", "clone", cloneURL, cloneDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		log.Error("failed to clone repository", "error", err)
		return
	}

	log.Info("Creating and configuring pulumi stack")

	os.Setenv("PULUMI_CONFIG_PASSPHRASE", "")

	stackCmd := exec.Command("pulumi", "stack", "select", "--create", minionName)
	stackCmd.Dir = cloneDir
	stackCmd.Stdout = os.Stdout
	stackCmd.Stderr = os.Stderr
	if err := stackCmd.Run(); err != nil {
		log.Error("failed to create or select pulumi stack", "error", err)
		return
	}

	configs := map[string]string{
		"salt-home:new_minion_name":  minionName,
		"salt-home:master_host":      c.masterHost,
		"salt-home:salt_minion_host": c.minionTarget,
	}

	for key, val := range configs {
		configCmd := exec.Command("pulumi", "config", "set", key, val)
		configCmd.Dir = cloneDir
		configCmd.Stdout = os.Stdout
		configCmd.Stderr = os.Stderr
		if err := configCmd.Run(); err != nil {
			log.Error("failed to set config", "key", key, "error", err)
			return
		}
	}

	grainsMap := make(map[string]interface{})
	for _, grain := range c.grains {
		parts := strings.SplitN(grain, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			if _, ok := grainsMap[key]; ok {
				if s, ok := grainsMap[key].([]string); ok {
					grainsMap[key] = append(s, value)
				} else {
					grainsMap[key] = []string{grainsMap[key].(string), value}
				}
			} else {
				grainsMap[key] = value
			}
		}
	}

	if len(grainsMap) > 0 {
		grainsJSON, err := json.Marshal(grainsMap)
		if err != nil {
			log.Error("failed to marshal grains to json", "error", err)
			return
		}

		configGrainsCmd := exec.Command("pulumi", "config", "set", "--path", "salt-home:grains", string(grainsJSON))
		configGrainsCmd.Dir = cloneDir
		configGrainsCmd.Stdout = os.Stdout
		configGrainsCmd.Stderr = os.Stderr
		if err := configGrainsCmd.Run(); err != nil {
			log.Error("failed to set grains config", "error", err)
			return
		}
	}

	log.Info("Running pulumi up")
	pulumiCmd := exec.Command("pulumi", "up", "--yes", "--skip-preview")
	pulumiCmd.Dir = cloneDir
	pulumiCmd.Stdout = os.Stdout
	pulumiCmd.Stderr = os.Stderr

	if err := pulumiCmd.Run(); err != nil {
		log.Error("failed to run pulumi up", "error", err)
		return
	}

	log.Info("Successfully added salt node")
}

func NewAddCommand() commands.BluePrintCommand {
	addCmd := &addCmd{}

	cobraCmd := &cobra.Command{
		Use:   "add <minion-name>",
		Short: "(Experimental) Adds a new salt node by cloning a git repository and running pulumi up.",
		Args:  cobra.ExactArgs(1),
		Run:   addCmd.run,
	}

	cobraCmd.Flags().StringVar(&addCmd.masterHost, "master-host", "", "The master host")
	cobraCmd.Flags().StringVar(&addCmd.minionTarget, "minion-target", "", "The minion target")
	cobraCmd.Flags().StringArrayVar(&addCmd.grains, "grain", []string{}, "Grains to set for the minion (e.g., roles=web, datacenter=nyc)")
	cobraCmd.MarkFlagRequired("master-host")
	cobraCmd.MarkFlagRequired("minion-target")

	addCmd.cobraCommand = cobraCmd
	return addCmd
}

func init() {
	commands.AddCommandToRegistry(NewAddCommand())
}
