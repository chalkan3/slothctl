package resources

import (
	"fmt"
	"os/exec"
	"strings"

	"slothctl/internal/log"
	"slothctl/pkg/bootstrap/common"
	"slothctl/pkg/statemanager"
)

// UserResource represents a system user to be managed.
type UserResource struct {
	Username string
	Password string // For initial creation/update, not stored in state
	UID      string // Desired UID
	GID      string // Desired GID
	Shell    string // Desired shell
}

// ID returns the unique identifier for the user resource.
func (u *UserResource) ID() string {
	return fmt.Sprintf("user:%s", u.Username)
}

// ReadCurrentState reads the current state of the user from the system.
func (u *UserResource) ReadCurrentState(dryRun bool) (map[string]interface{}, error) {
	log.Info("Reading current state for user", "username", u.Username, "dry_run", dryRun)

	// Check if user exists
	cmd := exec.Command("id", "-u", u.Username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "no such user") {
			return nil, nil // User does not exist
		}
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	// User exists, read details
	groupsCmd := exec.Command("id", "-Gn", u.Username)
	groupsOutput, err := groupsCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	currentGroups := strings.Fields(string(groupsOutput))
	isInRootGroup := false
	for _, group := range currentGroups {
		if group == "root" {
			isInRootGroup = true
			break
		}
	}

	return map[string]interface{}{
		"username": u.Username,
		"exists":   true,
		"inRootGroup": isInRootGroup,
	}, nil
}

// Diff compares the current state with the desired state and returns changes.
func (u *UserResource) Diff(currentState, desiredState map[string]interface{}) ([]statemanager.Change, error) {
	var changes []statemanager.Change

	// If user does not exist in current state, plan to create
	if currentState == nil || !currentState["exists"].(bool) {
		newValues := map[string]interface{}{"username": u.Username, "id": common.GenerateUUID()}
		if u.Password != "" {
			newValues["password"] = "[secret]" // Mask password
		}
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeCreate,
			ResourceID: u.ID(),
			NewValues:  newValues,
		})
		// Also plan to add to root group if creating
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeSetGroup,
			ResourceID: u.ID(),
			NewValues:  map[string]interface{}{"group": "root"},
		})
		return changes, nil
	}

	// Check if user needs to be added to root group
	if !currentState["inRootGroup"].(bool) {
		changes = append(changes, statemanager.Change{
			Type:       statemanager.ChangeTypeSetGroup,
			ResourceID: u.ID(),
			NewValues:  map[string]interface{}{"group": "root"},
		})
	}

	return changes, nil
}

// Apply applies the changes to the system.
func (u *UserResource) Apply(dryRun bool, changes []statemanager.Change) error {
	for _, change := range changes {
		log.Info("Applying change for user", "change_type", change.Type, "username", u.Username, "dry_run", dryRun)
		switch change.Type {
		case statemanager.ChangeTypeCreate:
			// Pass a goroutine name for CreateUser
			if err := common.CreateUser(common.GetRandomGoroutineName(), dryRun, u.Username, u.Password); err != nil {
				return fmt.Errorf("failed to create user %s: %w", u.Username, err)
			}
		case statemanager.ChangeTypeSetGroup:
			group := change.NewValues["group"].(string)
			if err := common.AddUserToGroup(common.GetRandomGoroutineName(), dryRun, u.Username, group); err != nil {
				return fmt.Errorf("failed to add user %s to group %s: %w", u.Username, group, err)
			}
		case statemanager.ChangeTypeUpdate:
			// Implement user update logic here
			log.Info("User update not yet implemented.", "username", u.Username)
		case statemanager.ChangeTypeDelete:
			// Implement user deletion logic here
			log.Info("User deletion not yet implemented.", "username", u.Username)
		case statemanager.ChangeTypeNoOp:
			log.Info("No operation for user", "username", u.Username)
		}
	}
	return nil
}
