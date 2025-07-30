package servermanager

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

const ( 
	ServerBucket = "servers"
	DefaultServerKey = "default_server"
)

// Server represents a managed server entry.
type Server struct {
	Name    string `json:"name"`
	Group   string `json:"group"`
	Context string `json:"context"`
	IP      string `json:"ip"`
	User    string `json:"user"`
	// Add other relevant fields like SSH key path, port, etc.
}

// Manager provides methods to interact with server data in BoltDB.
type Manager struct {
	db *bbolt.DB
}

// NewManager creates a new Server Manager instance.
func NewManager(db *bbolt.DB) *Manager {
	return &Manager{db: db}
}

// Init ensures the server bucket exists.
func (m *Manager) Init() error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ServerBucket))
		return err
	})
}

// SaveServer saves a server entry to the database.
func (m *Manager) SaveServer(server Server) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", ServerBucket)
		}

		key := []byte(fmt.Sprintf("%s:%s:%s", server.Group, server.Context, server.Name))
		encoded, err := json.Marshal(server)
		if err != nil {
			return err
		}
		return b.Put(key, encoded)
	})
}

// GetServer retrieves a server entry by its unique identifier (group:context:name).
func (m *Manager) GetServer(group, context, name string) (*Server, error) {
	var server Server
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", ServerBucket)
		}
		key := []byte(fmt.Sprintf("%s:%s:%s", group, context, name))
		val := b.Get(key)
		if val == nil {
			return fmt.Errorf("server %s:%s:%s not found", group, context, name)
		}
		return json.Unmarshal(val, &server)
	})
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// DeleteServer removes a server entry from the database.
func (m *Manager) DeleteServer(group, context, name string) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", ServerBucket)
		}
		key := []byte(fmt.Sprintf("%s:%s:%s", group, context, name))
		return b.Delete(key)
	})
}

// ListServers lists all server entries.
func (m *Manager) ListServers() ([]Server, error) {
	var servers []Server
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return nil // No servers yet
		}
		return b.ForEach(func(k, v []byte) error {
			var server Server
			if err := json.Unmarshal(v, &server); err != nil {
				return err
			}
			servers = append(servers, server)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// SetDefaultServer sets a server as the default.
func (m *Manager) SetDefaultServer(group, context, name string) error {
	return m.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", ServerBucket)
		}
		// Store the default server identifier
		key := []byte(DefaultServerKey)
		val := []byte(fmt.Sprintf("%s:%s:%s", group, context, name))
		return b.Put(key, val)
	})
}

// GetDefaultServer retrieves the default server.
func (m *Manager) GetDefaultServer() (string, string, string, error) {
	var group, context, name string
	err := m.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(ServerBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", ServerBucket)
		}
		val := b.Get([]byte(DefaultServerKey))
		if val == nil {
			return fmt.Errorf("no default server set")
		}
		parts := strings.Split(string(val), ":")
		if len(parts) != 3 {
			return fmt.Errorf("invalid default server format: %s", string(val))
		}
		group, context, name = parts[0], parts[1], parts[2]
		return nil
	})
	return group, context, name, err
}
