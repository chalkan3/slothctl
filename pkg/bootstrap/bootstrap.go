package bootstrap

import (
	"fmt"
	"sync" // For WaitGroup

	"slothctl/internal/log"
	"slothctl/pkg/bootstrap/incus"
	"slothctl/pkg/bootstrap/pass"
	"slothctl/pkg/bootstrap/salt"
	"slothctl/pkg/bootstrap/vault"
	"slothctl/pkg/bootstrap/common"
)

// RunControlPlaneBootstrap orchestrates the installation and configuration
// of SaltStack (master/minion), HashiCorp Vault, and Incus for a control plane.
func RunControlPlaneBootstrap(dryRun bool, saltUserPassword string) error {
	mainGoroutineName := "lady-guica" // Main goroutine name
	log.Info(fmt.Sprintf("%s is starting control plane bootstrapping process... %s", mainGoroutineName, log.GetRandomSlothEmoji()), "dry_run", dryRun)

	var wg sync.WaitGroup
	errChan := make(chan error, 4) // Buffer for 4 potential errors

	// 1. Install and Configure Vault (can run in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		goroutineName := common.GetRandomGoroutineName()
		log.Info(fmt.Sprintf("%s is starting Vault setup %s", goroutineName, log.GetRandomSlothEmoji()))
		if err := vault.InstallAndConfigureVault(goroutineName, dryRun); err != nil {
			errChan <- fmt.Errorf("%s: Vault bootstrap failed: %w", goroutineName, err)
		}
		log.Info(fmt.Sprintf("%s: Vault setup complete %s", goroutineName, log.GetRandomSlothEmoji()))
	}()

	// 2. Install and Configure Incus (can run in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		goroutineName := common.GetRandomGoroutineName()
		log.Info(fmt.Sprintf("%s is starting Incus setup %s", goroutineName, log.GetRandomSlothEmoji()))
		if err := incus.InstallAndConfigureIncus(goroutineName, dryRun); err != nil {
			errChan <- fmt.Errorf("%s: Incus bootstrap failed: %w", goroutineName, err)
		}
		log.Info(fmt.Sprintf("%s: Incus setup complete %s", goroutineName, log.GetRandomSlothEmoji()))
	}()

	// 3. Install and Configure SaltStack (master and minion) (can run in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		goroutineName := common.GetRandomGoroutineName()
		log.Info(fmt.Sprintf("%s is starting SaltStack setup %s", goroutineName, log.GetRandomSlothEmoji()))
		if err := salt.InstallAndConfigureSalt(goroutineName, dryRun, true, saltUserPassword); err != nil {
			errChan <- fmt.Errorf("%s: SaltStack bootstrap failed: %w", goroutineName, err)
		}
		log.Info(fmt.Sprintf("%s: SaltStack setup complete %s", goroutineName, log.GetRandomSlothEmoji()))
	}()

	// 4. Install and Configure Pass (can run in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		goroutineName := common.GetRandomGoroutineName()
		log.Info(fmt.Sprintf("%s is starting Pass setup %s", goroutineName, log.GetRandomSlothEmoji()))
		if err := pass.InstallAndConfigurePass(goroutineName, dryRun); err != nil {
			errChan <- fmt.Errorf("%s: Pass bootstrap failed: %w", goroutineName, err)
		}
		log.Info(fmt.Sprintf("%s: Pass setup complete %s", goroutineName, log.GetRandomSlothEmoji()))
	}()

	wg.Wait() // Wait for all goroutines to finish
	close(errChan) // Close the channel after all goroutines are done

	// Check for errors from goroutines
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	log.Info(fmt.Sprintf("%s: Control plane bootstrapping process complete. %s", mainGoroutineName, log.GetRandomSlothEmoji()))
	return nil
}