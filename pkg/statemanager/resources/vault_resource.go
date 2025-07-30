package resources

import (
	"fmt"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/statemanager"
	"github.com/chalkan3/slothctl/pkg/bootstrap/vault"
)

// VaultResource represents a HashiCorp Vault instance.
type VaultResource struct {
	ResourceID   string
	Name string
	// Add more Vault-specific attributes here (e.g., version, config path, address)
}

// ID returns the unique identifier for the Vault resource.
func (v *VaultResource) ID() string {
	return fmt.Sprintf("vault:%s", v.Name)
}

// ReadCurrentState reads the current state of the Vault instance from the system.
func (v *VaultResource) ReadCurrentState(dryRun bool) (map[string]interface{}, error) {
	log.Info("Reading current state for Vault", "name", v.Name, "dry_run", dryRun)
	// Simplified: In a real scenario, you'd check if Vault is installed, running, and get its config.
	// For now, assume it doesn't exist if we are planning to create it.
	return nil, nil // Assume it doesn't exist for now to always plan creation
}

// Diff compares the current state with the desired state and returns changes.
func (v *VaultResource) Diff(currentState, desiredState map[string]interface{}) ([]statemanager.Change, error) {
	var changes []statemanager.Change

	if currentState == nil {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeCreate,
			ResourceID: v.ID(),
			NewValues:  map[string]interface{}{"name": v.Name, "id": v.ResourceID},
		})
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeConfigure,
			ResourceID: v.ID(),
			NewValues:  map[string]interface{}{"address": "0.0.0.0:8200", "ui_enabled": true, "kind": "hashcorp-vault", "id": v.ResourceID},
		})
		return changes, nil
	}

	// For simplicity, assume only 'name' can be updated for now.
	// In a real scenario, you'd compare all relevant attributes.
	currentName := currentState["name"].(string)
	desiredName := desiredState["name"].(string)

	if currentName != desiredName {
		changes = append(changes, statemanager.Change{
			Type:        statemanager.ChangeTypeUpdate,
			ResourceID:  v.ID(),
			OldValues:   map[string]interface{}{"name": currentName},
			NewValues:   map[string]interface{}{"name": desiredName},
			DiffProperties: map[string]interface{}{"name": fmt.Sprintf("%s -> %s", currentName, desiredName)},
		})
	}

	if len(changes) == 0 {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeNoOp,
			ResourceID: v.ID(),
			Details:    map[string]interface{}{"message": "No changes detected"},
		})
	}

	return changes, nil
}

// Apply applies the changes to the system.
func (v *VaultResource) Apply(dryRun bool, changes []statemanager.Change) error {
	for _, change := range changes {
		log.Info("Applying change for Vault", "change_type", change.Type, "name", v.Name, "dry_run", dryRun)
		switch change.Type {
		case statemanager.ChangeTypeCreate:
			// Call the actual Vault installation/configuration logic here
			// For now, just log that it would be installed
			if err := vault.InstallAndConfigureVault(v.Name, dryRun); err != nil {
				return fmt.Errorf("failed to install and configure Vault: %w", err)
			}
		case statemanager.ChangeTypeConfigure:
			log.Info("Vault configuration would be performed here.", "dry_run", dryRun, "attributes", change.NewValues)
		case statemanager.ChangeTypeUpdate:
			log.Info("Vault update not yet implemented.", "dry_run", dryRun)
		case statemanager.ChangeTypeDelete:
			log.Info("Vault deletion not yet implemented.", "dry_run", dryRun)
		}
	}
	return nil
}
