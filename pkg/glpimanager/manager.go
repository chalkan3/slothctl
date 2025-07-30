package glpimanager

import (
	"encoding/json"
	"fmt"

	"github.com/chalkan3/slothctl/pkg/glpi"
	"go.etcd.io/bbolt"
)

const (
	GLPIBucket         = "glpi_instances"
	DefaultGLPIKey     = "default_glpi_instance"
	DefaultTicketIDKey = "default_ticket_id"
)

// Manager provides methods to interact with GLPI instance data in BoltDB.
type Manager struct {
	db *bbolt.DB
}

// NewManager creates a new GLPI Manager instance.
func NewManager(db *bbolt.DB) *Manager {
	return &Manager{db: db}
}

// Init ensures the GLPI bucket exists.
func (m *Manager) Init() error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(GLPIBucket))
		return err
	})
}

// SaveGLPIInstance saves a GLPI instance entry to the database.
func (m *Manager) SaveGLPIInstance(instance glpi.GLPIInstance) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}

		key := []byte(instance.Name)
		encoded, err := json.Marshal(instance)
		if err != nil {
			return err
		}
		return b.Put(key, encoded)
	})
}

// GetGLPIInstance retrieves a GLPI instance entry by its name.
func (m *Manager) GetGLPIInstance(name string) (*glpi.GLPIInstance, error) {
	var instance glpi.GLPIInstance
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		val := b.Get([]byte(name))
		if val == nil {
			return fmt.Errorf("GLPI instance %s not found", name)
		}
		return json.Unmarshal(val, &instance)
	})
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// DeleteGLPIInstance removes a GLPI instance entry from the database.
func (m *Manager) DeleteGLPIInstance(name string) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		return b.Delete([]byte(name))
	})
}

// ListGLPIInstances lists all GLPI instance entries.
func (m *Manager) ListGLPIInstances() ([]glpi.GLPIInstance, error) {
	var instances []glpi.GLPIInstance
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return nil // No instances yet
		}
		return b.ForEach(func(k, v []byte) error {
			var instance glpi.GLPIInstance
			if err := json.Unmarshal(v, &instance); err != nil {
				return err
			}
			instances = append(instances, instance)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// SetDefaultGLPIInstance sets a GLPI instance as the default.
func (m *Manager) SetDefaultGLPIInstance(name string) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		// Store the default GLPI instance name
		key := []byte(DefaultGLPIKey)
		val := []byte(name)
		return b.Put(key, val)
	})
}

// GetDefaultGLPIInstance retrieves the default GLPI instance name.
func (m *Manager) GetDefaultGLPIInstance() (string, error) {
	var name string
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		val := b.Get([]byte(DefaultGLPIKey))
		if val == nil {
			return fmt.Errorf("no default GLPI instance set")
		}
		name = string(val)
		return nil
	})
	return name, err
}

// GetGLPIClientForInstance retrieves a GLPI client for a given instance name.
func (m *Manager) GetGLPIClientForInstance(instanceName string) (*glpi.GLPIClient, error) {
	instance, err := m.GetGLPIInstance(instanceName)
	if err != nil {
		return nil, err
	}

	client := glpi.NewGLPIClient(instance.URL, instance.AppToken)

	// Authenticate using the stored user and password
	if err := client.Authenticate(instance.User, instance.Password); err != nil {
		return nil, fmt.Errorf("failed to authenticate with GLPI instance %s: %w", instanceName, err)
	}

	return client, nil
}

// GetDefaultGLPIClient retrieves a GLPI client for the default instance.
func (m *Manager) GetDefaultGLPIClient() (*glpi.GLPIClient, error) {
	defaultInstanceName, err := m.GetDefaultGLPIInstance()
	if err != nil {
		return nil, err
	}
	return m.GetGLPIClientForInstance(defaultInstanceName)
}

// SetDefaultTicketID sets a ticket ID as the default.
func (m *Manager) SetDefaultTicketID(ticketID int) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		key := []byte(DefaultTicketIDKey)
		val := []byte(fmt.Sprintf("%d", ticketID))
		return b.Put(key, val)
	})
}

// GetDefaultTicketID retrieves the default ticket ID.
func (m *Manager) GetDefaultTicketID() (int, error) {
	var ticketID int
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GLPIBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", GLPIBucket)
		}
		val := b.Get([]byte(DefaultTicketIDKey))
		if val == nil {
			return fmt.Errorf("no default ticket ID set")
		}
		_, err := fmt.Sscanf(string(val), "%d", &ticketID)
		if err != nil {
			return fmt.Errorf("failed to parse default ticket ID: %w", err)
		}
		return nil
	})
	return ticketID, err
}
