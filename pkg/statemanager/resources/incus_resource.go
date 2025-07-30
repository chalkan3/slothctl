package resources

import (
	"fmt"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/bootstrap/incus"
	"github.com/chalkan3/slothctl/pkg/statemanager"
)

// IncusResource represents an Incus host.
type IncusResource struct {
	ResourceID string
	Name       string
	// Add more Incus-specific attributes here (e.g., version, storage pools, networks)
}

// ID returns the unique identifier for the Incus resource.
func (i *IncusResource) ID() string {
	return fmt.Sprintf("incus:%s", i.Name)
}

// ReadCurrentState reads the current state of the Incus host from the system.
func (i *IncusResource) ReadCurrentState(dryRun bool) (map[string]interface{}, error) {
	log.Info("Reading current state for Incus", "name", i.Name, "dry_run", dryRun)
	// Simplified: In a real scenario, you'd check if Incus is installed and initialized.
	return nil, nil // Assume it doesn't exist for now to always plan creation
}

// Diff compares the current state with the desired state and returns changes.
func (i *IncusResource) Diff(currentState, desiredState map[string]interface{}) ([]statemanager.Change, error) {
	var changes []statemanager.Change

	if currentState == nil {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeCreate,
			ResourceID: i.ID(),
			NewValues:  map[string]interface{}{"name": i.Name, "id": i.ResourceID, "kind": "incus", "version": "12.12"},
		})
		return changes, nil
	}

	// For simplicity, assume only 'name' can be updated for now.
	// In a real scenario, you'd compare all relevant attributes.
	currentName := currentState["name"].(string)
	desiredName := desiredState["name"].(string)

	if currentName != desiredName {
		changes = append(changes, statemanager.Change{
			Type:           statemanager.ChangeTypeUpdate,
			ResourceID:     i.ID(),
			OldValues:      map[string]interface{}{"name": currentName},
			NewValues:      map[string]interface{}{"name": desiredName},
			DiffProperties: map[string]interface{}{"name": fmt.Sprintf("%s -> %s", currentName, desiredName)},
		})
	}

	if len(changes) == 0 {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeNoOp,
			ResourceID: i.ID(),
			Details:    map[string]interface{}{"message": "No changes detected"},
		})
	}

	return changes, nil
}

// Apply applies the changes to the system.
func (i *IncusResource) Apply(dryRun bool, changes []statemanager.Change) error {
	for _, change := range changes {
		log.Info("Applying change for Incus", "change_type", change.Type, "name", i.Name, "dry_run", dryRun)
		switch change.Type {
		case statemanager.ChangeTypeCreate:
			// Call the actual Incus installation/configuration logic here
			// For now, just log that it would be installed
			if err := incus.InstallAndConfigureIncus(i.Name, dryRun); err != nil {
				return fmt.Errorf("failed to install and configure Incus: %w", err)
			}
		case statemanager.ChangeTypeUpdate:
			log.Info("Incus update not yet implemented.", "dry_run", dryRun)
		case statemanager.ChangeTypeDelete:
			log.Info("Incus deletion not yet implemented.", "dry_run", dryRun)
		}
	}
	return nil
}
