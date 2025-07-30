package statemanager

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/chalkan3/slothctl/internal/log"
	"go.etcd.io/bbolt"
)

// Resource is the interface that all managed resources must implement.
type Resource interface {
	ID() string
	ReadCurrentState(dryRun bool) (map[string]interface{}, error)
	Diff(currentState, desiredState map[string]interface{}) ([]Change, error)
	Apply(dryRun bool, changes []Change) error
}

// ChangeType defines the type of change to a resource.
type ChangeType string

const (
	ChangeTypeCreate    ChangeType = "create"
	ChangeTypeUpdate    ChangeType = "update"
	ChangeTypeDelete    ChangeType = "delete"
	ChangeTypeNoOp      ChangeType = "no-op"
	ChangeTypeSetGroup  ChangeType = "set-group"
	ChangeTypeConfigure ChangeType = "configure"
)

// Change represents a planned or applied modification to a resource.
type Change struct {
	Type           ChangeType             `json:"type"`
	ResourceID     string                 `json:"resource_id"`
	NewValues      map[string]interface{} `json:"new_values,omitempty"`      // For create and update
	OldValues      map[string]interface{} `json:"old_values,omitempty"`      // For update and delete
	DiffProperties map[string]interface{} `json:"diff_properties,omitempty"` // Properties that changed
	Details        map[string]interface{} `json:"details,omitempty"`         // General details about the change
}

// StateManager manages the desired and current state of resources.
type StateManager struct {
	db     *bbolt.DB
	dryRun bool
}

// NewStateManager creates a new StateManager instance.
func NewStateManager(db *bbolt.DB, dryRun bool) *StateManager {
	return &StateManager{db: db, dryRun: dryRun}
}

// ReadState reads the current state of a resource from the BoltDB.
func (sm *StateManager) ReadState(resourceID string) (map[string]interface{}, error) {
	var state map[string]interface{}
	err := sm.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("slothctl_state"))
		if b == nil {
			return nil // Bucket doesn't exist yet, no state saved
		}
		data := b.Get([]byte(resourceID))
		if data == nil {
			return nil // No state for this resource ID
		}
		return json.Unmarshal(data, &state)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read state for %s: %w", resourceID, err)
	}
	return state, nil
}

// WriteState writes the current state of a resource to the BoltDB.
func (sm *StateManager) WriteState(resourceID string, state map[string]interface{}) error {
	return sm.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("slothctl_state"))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		data, err := json.Marshal(state)
		if err != nil {
			return fmt.Errorf("marshal state: %w", err)
		}
		return b.Put([]byte(resourceID), data)
	})
}

// Plan compares the desired state with the current state and generates a plan of changes.
func (sm *StateManager) Plan(desiredResources []Resource) ([]Change, error) {
	log.Info("Generating execution plan...")
	var allChanges []Change

	for _, desiredRes := range desiredResources {
		resourceID := desiredRes.ID()
		log.Info("Planning for resource", "id", resourceID, "type", reflect.TypeOf(desiredRes).Elem().Name())

		// Read current state from system
		currentState, err := desiredRes.ReadCurrentState(sm.dryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to read current state for %s: %w", resourceID, err)
		}

		// For simplicity, we'll use desiredState as the source of truth for diffing against current.
		// A more complex state manager would diff desired vs. lastKnown, and then current vs. desired.
		desiredStateForDiff, err := desiredRes.ReadCurrentState(true) // Read desired state for diffing (no side effects)
		if err != nil {
			return nil, fmt.Errorf("failed to read desired state for diffing for %s: %w", resourceID, err)
		}

		changes, err := desiredRes.Diff(currentState, desiredStateForDiff)
		if err != nil {
			return nil, fmt.Errorf("failed to diff resource %s: %w", resourceID, err)
		}

		if len(changes) == 0 {
			log.Info("No changes detected for resource", "id", resourceID)
			continue
		}

		log.Info("Changes planned for resource", "id", resourceID, "changes_count", len(changes))
		allChanges = append(allChanges, changes...)
	}

	log.Info("Execution plan generated.", "total_changes", len(allChanges))
	return allChanges, nil
}

// Apply applies the planned changes to the system and updates the state in BoltDB.
func (sm *StateManager) Apply(changes []Change, desiredResources []Resource) error {
	log.Info("Applying changes...", "total_changes", len(changes), "dry_run", sm.dryRun)

	resourceMap := make(map[string]Resource)
	for _, res := range desiredResources {
		resourceMap[res.ID()] = res
	}

	for _, change := range changes {
		log.Info("Applying change", "type", change.Type, "resource_id", change.ResourceID, "details", change.Details, "dry_run", sm.dryRun)

		res, ok := resourceMap[change.ResourceID]
		if !ok {
			log.Error("Resource not found for change", "resource_id", change.ResourceID)
			continue
		}

		if err := res.Apply(sm.dryRun, []Change{change}); err != nil {
			return fmt.Errorf("failed to apply change for %s: %w", change.ResourceID, err)
		}

		// After applying, read the new current state and save it to BoltDB
		if !sm.dryRun {
			log.Info("Updating state in DB for resource", "id", change.ResourceID)
			updatedState, err := res.ReadCurrentState(sm.dryRun) // Read actual state after apply
			if err != nil {
				log.Error("Failed to read updated state after apply", "resource_id", change.ResourceID, "error", err)
				// Continue, but log the error
			}
			if updatedState != nil {
				if err := sm.WriteState(change.ResourceID, updatedState); err != nil {
					log.Error("Failed to write updated state to DB", "resource_id", change.ResourceID, "error", err)
					// Continue, but log the error
				}
			}
		}
	}

	log.Info("Changes applied.", "dry_run", sm.dryRun)
	return nil
}
