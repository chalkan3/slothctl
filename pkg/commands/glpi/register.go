package glpi

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/chalkan3/slothctl/internal/log"
	"github.com/chalkan3/slothctl/pkg/commands"
	"github.com/chalkan3/slothctl/pkg/config"
	"github.com/chalkan3/slothctl/pkg/glpi"
	"github.com/chalkan3/slothctl/pkg/glpimanager"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"golang.org/x/term"
)

// registerCmd represents the 'glpi register' command
type registerCmd struct{}

func (c *registerCmd) Parent() string {
	return "glpi"
}

func (c *registerCmd) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Short: "Registers a new GLPI instance",
		Long:  `Registers a new GLPI instance with its URL, App Token, and user details.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			url, _ := cmd.Flags().GetString("url")
			appToken, _ := cmd.Flags().GetString("app-token")
			user, _ := cmd.Flags().GetString("user")

			if url == "" || appToken == "" || user == "" {
				return fmt.Errorf("url, app-token, and user flags are required")
			}

			fmt.Print("Enter GLPI password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			password := string(bytePassword)
			fmt.Println("") // Newline after password input

			// Initialize BoltDB
			dbPath := os.ExpandEnv(config.AppConfig.DatabasePath)
			db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return fmt.Errorf("failed to open BoltDB: %w", err)
			}
			defer db.Close()

			sm := glpimanager.NewManager(db)
			if err := sm.Init(); err != nil {
				return fmt.Errorf("failed to initialize GLPI manager: %w", err)
			}

			instance := glpi.GLPIInstance{
				Name:     name,
				URL:      url,
				AppToken: appToken,
				User:     user,
				Password: password,
			}

			// Test authentication before saving
			client := glpi.NewGLPIClient(instance.URL, instance.AppToken)
			if err := client.Authenticate(instance.User, instance.Password); err != nil {
				return fmt.Errorf("authentication failed for GLPI instance %s: %w", name, err)
			}
			log.Info("Authentication successful. Saving GLPI instance...")

			if err := sm.SaveGLPIInstance(instance); err != nil {
				return fmt.Errorf("failed to save GLPI instance: %w", err)
			}

			log.Info("GLPI instance registered successfully.", "name", name)
			return nil
		},
	}

	cmd.Flags().StringP("url", "u", "", "GLPI instance URL (e.g., https://your-glpi.com/glpi)")
	cmd.Flags().StringP("app-token", "a", "", "GLPI App Token")
	cmd.Flags().StringP("user", "l", "", "GLPI username")

	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("app-token")
	cmd.MarkFlagRequired("user")

	return cmd
}

func init() {
	commands.AddCommandToRegistry(&registerCmd{})
}
