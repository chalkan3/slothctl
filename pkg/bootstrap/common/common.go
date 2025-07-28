package common

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"math/rand" // For random goroutine names
	"time"      // For seeding rand

	"github.com/google/uuid"
	"slothctl/internal/log"
)

// RunCommand executes a shell command and logs its output.
func RunCommand(goroutineName string, dryRun bool, stdin io.Reader, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = log.NewWriter(log.Info)
	cmd.Stderr = log.NewWriter(log.Error)
	cmd.Stdin = stdin

	log.Info(fmt.Sprintf("%s is handling command: %s", goroutineName, cmd.String()), "dry_run", dryRun)

	if dryRun {
		log.Info(fmt.Sprintf("%s: Dry run: Command not executed.", goroutineName))
		return nil
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %s %w", cmd.String(), err)
	}
	return nil
}

// InstallPackages installs a list of packages using pacman.
func InstallPackages(goroutineName string, dryRun bool, packages []string) error {
	log.Info(fmt.Sprintf("%s is installing packages: %v", goroutineName, packages), "dry_run", dryRun)
	args := []string{"--noconfirm", "-S"}
	args = append(args, packages...)
	return RunCommand(goroutineName, dryRun, nil, "sudo", append([]string{"pacman"}, args...)...)
}

// CreateUser creates a system user with a specified password.
func CreateUser(goroutineName string, dryRun bool, username, password string) error {
	log.Info(fmt.Sprintf("%s is creating system user: %s", goroutineName, username), "dry_run", dryRun)

	// Check if user already exists
	if err := RunCommand(goroutineName, dryRun, nil, "id", "-u", username); err == nil {
		log.Info(fmt.Sprintf("%s: User already exists.", goroutineName), "username", username)
		return nil
	}

	// Create user
	if err := RunCommand(goroutineName, dryRun, nil, "sudo", "useradd", "-m", "-s", "/bin/bash", username); err != nil {
		return fmt.Errorf("failed to create user %s: %w", username, err)
	}

	// Set password
	if password != "" {
		log.Info(fmt.Sprintf("%s is setting password for user: %s", goroutineName, username), "dry_run", dryRun)
		cmd := exec.Command("sudo", "chpasswd")
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%s:%s", username, password))
		cmd.Stdout = log.NewWriter(log.Info)
		cmd.Stderr = log.NewWriter(log.Error)

		log.Info(fmt.Sprintf("%s is running command: %s", goroutineName, cmd.String()), "dry_run", dryRun)

		if dryRun {
			log.Info(fmt.Sprintf("%s: Dry run: Command not executed.", goroutineName))
			return nil
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed: %s %w", cmd.String(), err)
		}
	}

	log.Info(fmt.Sprintf("%s: System user created successfully.", goroutineName), "username", username)
	return nil
}

// GenerateUUID generates a new UUID string.
func GenerateUUID() string {
	return uuid.New().String()
}

var goroutineNames = []string{"maria-guica", "keity-guica", "mini-guica"}

// GetRandomGoroutineName returns a random goroutine name.
func GetRandomGoroutineName() string {
	rand.Seed(time.Now().UnixNano())
	return goroutineNames[rand.Intn(len(goroutineNames))]
}

// AddUserToGroup adds a user to a specified group.
func AddUserToGroup(goroutineName string, dryRun bool, username, group string) error {
	log.Info(fmt.Sprintf("%s is adding user %s to group %s", goroutineName, username, group), "dry_run", dryRun)
	return RunCommand(goroutineName, dryRun, nil, "sudo", "usermod", "-aG", group, username)
}