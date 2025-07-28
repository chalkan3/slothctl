package resources

import (
	"fmt"

	"slothctl/internal/log"
	"slothctl/pkg/statemanager"
	"slothctl/pkg/bootstrap/salt"
)

// SaltMinionResource represents a Salt Minion instance.
type SaltMinionResource struct {
	ResourceID   string
	Name string
	// Add more Salt Minion-specific attributes here (e.g., master address, minion ID)
}

// ID returns the unique identifier for the Salt Minion resource.
func (s *SaltMinionResource) ID() string {
	return fmt.Sprintf("salt_minion:%s", s.Name)
}

// ReadCurrentState reads the current state of the Salt Minion from the system.
func (s *SaltMinionResource) ReadCurrentState(dryRun bool) (map[string]interface{}, error) {
	log.Info("Reading current state for Salt Minion", "name", s.Name, "dry_run", dryRun)
	// Simplified: In a real scenario, you'd check if Salt Minion is installed and running.
	return nil, nil // Assume it doesn't exist for now to always plan creation
}

// Diff compares the current state with the desired state and returns changes.
func (s *SaltMinionResource) Diff(currentState, desiredState map[string]interface{}) ([]statemanager.Change, error) {
	var changes []statemanager.Change

	if currentState == nil {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeCreate,
			ResourceID: s.ID(),
			NewValues:  map[string]interface{}{"name": s.Name, "id": s.ResourceID, "kind": "salt_minion", "version": "12.12"},
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
			ResourceID:  s.ID(),
			OldValues:   map[string]interface{}{"name": currentName},
			NewValues:   map[string]interface{}{"name": desiredName},
			DiffProperties: map[string]interface{}{"name": fmt.Sprintf("%s -> %s", currentName, desiredName)},
		})
	}

	if len(changes) == 0 {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeNoOp,
			ResourceID: s.ID(),
			Details:    map[string]interface{}{"message": "No changes detected"},
		})
	}

	return changes, nil
}

// Apply applies the changes to the system.
func (s *SaltMinionResource) Apply(dryRun bool, changes []statemanager.Change) error {
	for _, change := range changes {
		log.Info("Applying change for Salt Minion", "change_type", change.Type, "name", s.Name, "dry_run", dryRun)
		switch change.Type {
		case statemanager.ChangeTypeCreate:
			// Call the actual Salt Minion installation/configuration logic here
			// For now, just log that it would be installed
			if err := salt.InstallAndConfigureSalt(s.Name, dryRun, false, ""); err != nil {
				return fmt.Errorf("failed to install and configure Salt Minion: %w", err)
			}
		case statemanager.ChangeTypeUpdate:
			log.Info("Salt Minion update not yet implemented.", "dry_run", dryRun)
		case statemanager.ChangeTypeDelete:
			log.Info("Salt Minion deletion not yet implemented.", "dry_run", dryRun)
		}
	}
	return nil
}
