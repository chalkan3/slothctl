package salt

import (
	"fmt"

	"slothctl/internal/log"
	"slothctl/pkg/bootstrap/common"
)

const (
	saltMasterConfigPath = "/etc/salt/master"
	saltMinionConfigPath = "/etc/salt/minion"
	gitFsRepoURL         = "https://github.com/your-org/your-salt-states.git" // Mock URL
	saltUserName         = "saltuser"
)

// InstallAndConfigureSalt installs and configures SaltStack (master and/or minion).
func InstallAndConfigureSalt(goroutineName string, dryRun bool, isMaster bool, saltUserPassword string) error {
	log.Info(fmt.Sprintf("%s is starting SaltStack installation and configuration...", goroutineName), "dry_run", dryRun)

	packages := []string{"salt"}
	if isMaster {
		packages = append(packages, "salt-master")
	}

	// Install Salt packages
	// common.InstallPackages already handles dryRun and sudo
	if err := common.InstallPackages(goroutineName, dryRun, packages); err != nil {
		return fmt.Errorf("failed to install Salt packages: %w", err)
	}

	// Create dedicated Salt user if password is provided
	if saltUserPassword != "" {
		log.Info(fmt.Sprintf("%s is creating dedicated Salt user...", goroutineName), "username", saltUserName, "dry_run", dryRun)
		if err := common.CreateUser(goroutineName, dryRun, saltUserName, saltUserPassword); err != nil {
			return fmt.Errorf("failed to create Salt user: %w", err)
		}
		log.Info(fmt.Sprintf("%s: Dedicated Salt user created.", goroutineName))
	}

	// Configure Salt Master (if applicable)
	if isMaster {
		log.Info(fmt.Sprintf("%s is configuring Salt Master...", goroutineName), "dry_run", dryRun)
		masterConfigContent := fmt.Sprintf(`
fileserver_backend:
  - roots
  - git

gitfs_remotes:
  - %s

# External authentication for Salt Master
external_auth:
  pam:
    %s:
      - .*

# ACL for the dedicated Salt user
client_acl:
  %s:
    - .*

`, gitFsRepoURL, saltUserName, saltUserName)

		// Use a here-document to write multi-line content to file
		cmdStr := fmt.Sprintf("cat <<EOF > %s\n%sEOF", saltMasterConfigPath, masterConfigContent)
		if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "sh", "-c", cmdStr); err != nil {
			return fmt.Errorf("failed to write Salt Master config: %w", err)
		}
		log.Info(fmt.Sprintf("%s: Salt Master configured.", goroutineName))

		// Enable and start salt-master service
		log.Info(fmt.Sprintf("%s is enabling and starting salt-master service...", goroutineName), "dry_run", dryRun)
		if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "enable", "salt-master"); err != nil {
			return fmt.Errorf("failed to enable salt-master service: %w", err)
		}
		if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "start", "salt-master"); err != nil {
			return fmt.Errorf("failed to start salt-master service: %w", err)
		}
		log.Info(fmt.Sprintf("%s: Salt Master service started.", goroutineName))
	}

	// Configure Salt Minion
	log.Info(fmt.Sprintf("%s is configuring Salt Minion...", goroutineName), "dry_run", dryRun)
	minionConfigContent := `
master: 127.0.0.1 # Assuming master is on the same machine for control-plane

# Optional: Minion ID (defaults to hostname)
# id: my-minion-id
`

	// Use a here-document to write multi-line content to file
	cmdStr := fmt.Sprintf("cat <<EOF > %s\n%sEOF", saltMinionConfigPath, minionConfigContent)
		if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "sh", "-c", cmdStr); err != nil {
			return fmt.Errorf("failed to write Salt Minion config: %w", err)
		}
	log.Info(fmt.Sprintf("%s: Salt Minion configured.", goroutineName))

	// Enable and start salt-minion service
	log.Info(fmt.Sprintf("%s is enabling and starting salt-minion service...", goroutineName), "dry_run", dryRun)
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "enable", "salt-minion"); err != nil {
		return fmt.Errorf("failed to enable salt-minion service: %w", err)
	}
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "start", "salt-minion"); err != nil {
		return fmt.Errorf("failed to start salt-minion service: %w", err)
	}
	log.Info(fmt.Sprintf("%s: Salt Minion service started.", goroutineName))

	log.Info(fmt.Sprintf("%s: SaltStack installation and configuration complete.", goroutineName))
	return nil
}