package configure

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"slothctl/internal/log"
	"slothctl/pkg/bootstrap"
	"slothctl/pkg/commands"
	"slothctl/pkg/config"
	"slothctl/pkg/statemanager"
	"slothctl/pkg/statemanager/resources"
	"slothctl/pkg/bootstrap/common" // Import common for GenerateUUID
)

// initCmd represents the 'configure init' command
type initCmd struct{}

func (c *initCmd) Parent() string {
	return "configure"
}

func (c *initCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes slothctl configuration and sets up the embedded database and control plane components",
		Long:  "The init command performs initial setup for slothctl, including creating default configuration files, setting up the embedded BoltDB database, and installing/configuring control plane components like Vault, Incus, and SaltStack.",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Initializing slothctl...")

			// Initialize default config file
			if err := config.InitConfig(); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			// Reload config to ensure default values are loaded into AppConfig
			if err := config.LoadConfig(); err != nil {
				return fmt.Errorf("failed to reload config after initialization: %w", err)
			}
			log.Info("Database path from config", "path", config.AppConfig.DatabasePath)

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			log.Info("Expanded database path", "path", dbPath)

			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			// Create a default bucket if it doesn't exist
			err = db.Update(func(tx *bbolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte("slothctl_data"))
				if err != nil {
					return fmt.Errorf("create bucket: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to create default database bucket: %w", err)
			}

			log.Info("Embedded database initialized successfully.")

			// Get flags
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			saltUserPassword, _ := cmd.Flags().GetString("salt-user-password")
			mode, _ := cmd.Flags().GetString("mode")
			username, _ := cmd.Flags().GetString("user")

			// Initialize StateManager
			sm := statemanager.NewStateManager(db, dryRun)

			// Define desired state (example: a system user)
			desiredResources := []statemanager.Resource{
				&resources.UserResource{Username: username, Password: "supersecret"},
				&resources.VaultResource{ResourceID: common.GenerateUUID(), Name: "sloth-vault"},
				&resources.IncusResource{ResourceID: common.GenerateUUID(), Name: "sloth-incus-host"},
				&resources.SaltMasterResource{ResourceID: common.GenerateUUID(), Name: "sloth-salt-master"},
				&resources.SaltMinionResource{ResourceID: common.GenerateUUID(), Name: "sloth-salt-minion"},
			}

			var changes []statemanager.Change
			var planErr error

			// Always generate a plan first
			changes, planErr = sm.Plan(desiredResources)
			if planErr != nil {
				return fmt.Errorf("failed to generate plan: %w", planErr)
			}

			// Print plan in Pulumi-like format
			if len(changes) == 0 {
				fmt.Println("No changes. Your infrastructure matches the configuration.")
			} else {
				fmt.Println("│     Type                     Attributes              Plan       Info                                                                                    │")

				var createdUserID string // To store the ID of the created user for @depends

				// Group Incus, Salt Master, and Salt Minion under Package
				packageChanges := make(map[string][]statemanager.Change)
				otherChanges := []statemanager.Change{}

				for _, change := range changes {
					resourceKind := strings.Split(change.ResourceID, ":")[0]
					switch resourceKind {
					case "incus", "salt_master", "salt_minion":
						packageChanges["Package"] = append(packageChanges["Package"], change)
					default:
						otherChanges = append(otherChanges, change)
					}
				}

				// Print other changes first
				for _, change := range otherChanges {
					resourceKind := strings.Split(change.ResourceID, ":")[0]
					_ = strings.Split(change.ResourceID, ":")[1] // resourceName is declared but not used in this block, so ignore it
					_ = strings.Split(change.ResourceID, ":")[1] // resourceName is declared but not used in this scope // Re-declare for use in the switch statement

					switch resourceKind {
					case "user":
						if change.Type == statemanager.ChangeTypeCreate {
							fmt.Printf("│            +   [Core:User]                                                                                                                     │\n")
							fmt.Printf("│        => create                                                                                                                          │\n")
							if change.NewValues != nil {
								for k, v := range change.NewValues {
									val := fmt.Sprintf("%v", v)
									if strings.Contains(k, "password") || strings.Contains(k, "secret") || strings.Contains(k, "token") {
										val = "[secret]"
									}
									if k == "id" {
                                        createdUserID = val // Capture the user ID
                                        val = "{algum_uuid}"
                                    }
                                    fmt.Printf("│                                 |─ %-14s        %s                                                        │\n", k, val)
								}
							}
						} else if change.Type == statemanager.ChangeTypeSetGroup {
							fmt.Printf("│                                                                                         +   [Core:Group]                                        │\n")
							fmt.Printf("│                                                                                                      +     => add-user-to-group          │\n")
							fmt.Printf("│                                                                                                                               @depends Core:User:create%%%s                                                                                        │\n", createdUserID) // Use the captured user ID
							if change.NewValues != nil {
								for k, v := range change.NewValues {
									val := fmt.Sprintf("%v", v)
									fmt.Printf("│                                                                                                 +        ├─   %-14s -> %s                                                                                   │\n", k, val)
								}
							}
						}
					case "vault":
						if change.Type == statemanager.ChangeTypeCreate || change.Type == statemanager.ChangeTypeConfigure {
							fmt.Printf("│                            +   [Core:SecretManager]                                                                           │\n")
							fmt.Printf("│                                => start-daemon                                                                                         │\n")
							if change.NewValues != nil {
								for k, v := range change.NewValues {
									val := fmt.Sprintf("%v", v)
									if strings.Contains(k, "password") || strings.Contains(k, "secret") || strings.Contains(k, "token") {
										val = "[secret]"
									}
									if k == "kind" {
										val = "hashcorp-vault"
									}
									if k == "id" {
										val = "{algum_uuid}"
									}
									fmt.Printf("│                                   +      ├─ %-14s -> %s                                                        │\n", k, val)
								}
							}
						}
					}
				}

				                // Print package changes
                if len(packageChanges) > 0 {
                    fmt.Printf("│                +   [Core:PackageManager]                                                                                                      │\n")
                    fmt.Printf("│                                    +   │   => install                                                                                         │\n")
                    fmt.Printf("│                                                        @depends None                                                             │\n")
                    for _, changes := range packageChanges {
                        for _, change := range changes {
                            _ = strings.Split(change.ResourceID, ":")[0] // resourceType
                            resourceName := strings.Split(change.ResourceID, ":")[1]

                            fmt.Printf("│                                                                            +   │        ├─ name       ->  %s                    │\n", resourceName)
                            if change.NewValues != nil {
                                for k, v := range change.NewValues {
                                    val := fmt.Sprintf("%v", v)
                                    if strings.Contains(k, "password") || strings.Contains(k, "secret") || strings.Contains(k, "token") {
                                        val = "[secret]"
                                    }
                                    if k == "id" {
                                        val = "{algum_uuid}"
                                    }
                                    fmt.Printf("│                                                                                                                    +   │        ├─ %-14s -> %s                                                                                 │\n", k, val)
                                }
                            }
                        }
                    }
                }

				fmt.Println("\nResources:")
				fmt.Printf("    + %d to create\n", countChanges(changes, statemanager.ChangeTypeCreate))
				fmt.Printf("    ~ %d to update\n", countChanges(changes, statemanager.ChangeTypeUpdate))
				fmt.Printf("    - %d to delete\n", countChanges(changes, statemanager.ChangeTypeDelete))
			}

			switch mode {
			case "plan":
				log.Info("Plan mode: Exiting after plan generation.")
				return nil
			case "apply":
				log.Info("Apply mode: Applying changes.")
				if dryRun {
					log.Info("Dry run enabled. Changes will not be applied.")
					return nil
				}
				if err := sm.Apply(changes, desiredResources); err != nil {
					return fmt.Errorf("failed to apply changes: %w", err)
				}
			case "bootstrap": // Default mode for initial setup
				log.Info("Bootstrap mode: Running control plane bootstrap.")
				if err := bootstrap.RunControlPlaneBootstrap(dryRun, saltUserPassword); err != nil {
					return fmt.Errorf("control plane bootstrap failed: %w", err)
				}
			default:
				return fmt.Errorf("invalid mode specified: %s. Use 'plan', 'apply', or 'bootstrap'.", mode)
			}

			log.Info("slothctl initialization complete.")
			return nil
		},
	}
	cmd.Flags().Bool("dry-run", false, "If true, commands will only be logged and not executed.")
	cmd.Flags().String("salt-user-password", "", "Password for the dedicated SaltStack user (required for advanced config).")
	cmd.Flags().String("mode", "bootstrap", "Operation mode: 'plan', 'apply', or 'bootstrap' (default). ")
	cmd.Flags().String("user", "myadmin", "Username for the primary system user.")
	return cmd
}

func init() {
	commands.AddCommandToRegistry(&initCmd{})
}

// Helper function to count changes by type
func countChanges(changes []statemanager.Change, changeType statemanager.ChangeType) int {
	count := 0
	for _, change := range changes {
		if change.Type == changeType {
			count++
		}
	}
	return count
}
