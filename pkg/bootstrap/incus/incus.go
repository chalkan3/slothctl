package incus

import (
	"fmt"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/bootstrap/common"
)

// InstallAndConfigureIncus installs and configures Incus.
func InstallAndConfigureIncus(goroutineName string, dryRun bool) error {
	log.Info(fmt.Sprintf("%s is starting Incus installation and configuration...", goroutineName), "dry_run", dryRun)

	// Install Incus package
	packages := []string{"incus"}
	if err := common.InstallPackages(goroutineName, dryRun, packages); err != nil {
		return fmt.Errorf("failed to install Incus package: %w", err)
	}

	// Initialize Incus (non-interactive)
	// This command sets up storage pools, networks, etc.
	log.Info(fmt.Sprintf("%s is initializing Incus (non-interactive)...", goroutineName), "dry_run", dryRun)
	// A basic non-interactive setup might look like this:
	// incus init --auto --storage-backend dir --network-address 0.0.0.0 --network-port 8443
	// A fully automated setup would require more specific flags or a preseed file.

	// Check if Incus is already initialized
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "incus", "info"); err == nil {
		log.Info(fmt.Sprintf("%s: Incus is already initialized.", goroutineName))
	} else {
		log.Info(fmt.Sprintf("%s is running incus init --auto...", goroutineName), "dry_run", dryRun)
		if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "incus", "init", "--auto"); err != nil {
			return fmt.Errorf("failed to initialize Incus: %w", err)
		}
		log.Info(fmt.Sprintf("%s: Incus initialized.", goroutineName))
	}

	log.Info(fmt.Sprintf("%s: Incus installation and configuration complete.", goroutineName))
	return nil
}
