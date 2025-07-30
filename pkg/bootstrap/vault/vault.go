package vault

import (
	"fmt"
	"path/filepath"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/bootstrap/common"
)

const (
	vaultConfigPath = "/etc/vault/vault.hcl"
	vaultDataPath   = "/opt/vault/data"
)

// InstallAndConfigureVault installs and configures HashiCorp Vault.
func InstallAndConfigureVault(goroutineName string, dryRun bool) error {
	log.Info(fmt.Sprintf("%s is starting HashiCorp Vault installation and configuration...", goroutineName), "dry_run", dryRun)

	// Install Vault package
	// Vault is typically distributed as a pre-compiled binary or via a specific repository.
	// For Arch Linux, it's usually in the community repository.
	packages := []string{"vault"}
	if err := common.InstallPackages(goroutineName, dryRun, packages); err != nil {
		return fmt.Errorf("failed to install Vault package: %w", err)
	}

	// Create Vault data directory
	log.Info(fmt.Sprintf("%s is creating Vault data directory...", goroutineName), "dry_run", dryRun)
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "mkdir", "-p", vaultDataPath); err != nil {
		return fmt.Errorf("failed to create Vault data directory: %w", err)
	}
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "chown", "vault:vault", vaultDataPath); err != nil {
		return fmt.Errorf("failed to set ownership for Vault data directory: %w", err)
	}

	// Configure Vault
	log.Info(fmt.Sprintf("%s is configuring Vault...", goroutineName), "dry_run", dryRun)
	vaultConfigContent := fmt.Sprintf(`
storage "file" {
  path = "%s"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

ui = true

`, vaultDataPath)

	// Ensure /etc/vault directory exists
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "mkdir", "-p", filepath.Dir(vaultConfigPath)); err != nil {
		return fmt.Errorf("failed to create /etc/vault directory: %w", err)
	}

	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "sh", "-c", fmt.Sprintf("echo \"%s\" > %s", vaultConfigContent, vaultConfigPath)); err != nil {
		return fmt.Errorf("failed to write Vault config: %w", err)
	}
	log.Info(fmt.Sprintf("%s: Vault configured.", goroutineName))

	// Enable and start vault service
	log.Info(fmt.Sprintf("%s is enabling and starting vault service...", goroutineName), "dry_run", dryRun)
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "enable", "vault"); err != nil {
		return fmt.Errorf("failed to enable vault service: %w", err)
	}
	if err := common.RunCommand(goroutineName, dryRun, nil, "sudo", "systemctl", "start", "vault"); err != nil {
		return fmt.Errorf("failed to start vault service: %w", err)
	}
	log.Info(fmt.Sprintf("%s: Vault service started.", goroutineName))

	// Note: Vault requires initialization and unsealing after startup.
	// This is typically a manual or automated process outside of basic installation.
	// For bootstrapping, you might consider using 'vault operator init' and 'vault operator unseal'
	// with a simple file backend for development/testing purposes.

	log.Info(fmt.Sprintf("%s: HashiCorp Vault installation and configuration complete.", goroutineName))
	return nil
}