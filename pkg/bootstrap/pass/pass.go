package pass

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"slothctl/internal/log"
	"slothctl/pkg/bootstrap/common"
)

// InstallAndConfigurePass installs and configures GNU Pass.
func InstallAndConfigurePass(goroutineName string, dryRun bool) error {
	log.Info(fmt.Sprintf("%s is starting GNU Pass installation and configuration...", goroutineName), "dry_run", dryRun)

	// Install pass and gnupg packages
	packages := []string{"pass", "gnupg"}
	if err := common.InstallPackages(goroutineName, dryRun, packages); err != nil {
		return fmt.Errorf("failed to install pass/gnupg packages: %w", err)
	}

	// Check for existing GPG key
	log.Info(fmt.Sprintf("%s is checking for existing GPG key...", goroutineName), "dry_run", dryRun)
	// This is a simplified check. A real check would parse `gpg --list-keys` output.
	// For dry-run, we just log the command.
	if err := common.RunCommand(goroutineName, dryRun, nil, "gpg", "--list-keys"); err != nil {
		log.Warn(fmt.Sprintf("%s: No GPG key found. You might need to generate one manually or automate it.", goroutineName))
		log.Warn(fmt.Sprintf("%s: Example: gpg --full-generate-key", goroutineName))
	}

	// Initialize pass repository
	passDir := filepath.Join(os.ExpandEnv("$HOME"), ".password-store")
	if _, err := os.Stat(passDir); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("%s is initializing pass repository...", goroutineName), "path", passDir, "dry_run", dryRun)
		// This command requires a GPG key to be present and selected.
		// For dry-run, we just log the command.
		var stdin io.Reader = nil
		if !dryRun {
			log.Info(fmt.Sprintf("%s: Please enter your GPG passphrase if prompted by pass init.", goroutineName))
			stdin = os.Stdin // Pass os.Stdin for interactive passphrase input
		}
		if err := common.RunCommand(goroutineName, dryRun, stdin, "pass", "init", "YOUR_GPG_KEY_ID"); err != nil {
			return fmt.Errorf("failed to initialize pass repository: %w", err)
		}
	} else {
		log.Info(fmt.Sprintf("%s: Pass repository already exists.", goroutineName), "path", passDir)
	}

	log.Info(fmt.Sprintf("%s: GNU Pass installation and configuration complete.", goroutineName))
	return nil
}