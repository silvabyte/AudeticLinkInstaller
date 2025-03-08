package system

import (
	"fmt"
	"os/exec"
)

// InstallDependencies installs all required system packages
func InstallDependencies() error {
	deps := []string{
		"git",
		"python3-pip",
		"python3-venv",
		"python3-gpiozero",
		"sox",
		"libsox-fmt-all",
		"netcat-openbsd",
		"rpi-connect-lite",
		"i2c-tools",
		"python3-smbus",
		"alsa-utils",
	}

	// Update package list
	if err := execCommand("apt", "update"); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Upgrade system
	if err := execCommand("apt", "full-upgrade", "-y"); err != nil {
		return fmt.Errorf("failed to upgrade system: %w", err)
	}

	// Install packages
	args := append([]string{"install", "-y"}, deps...)
	if err := execCommand("apt", args...); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil // Silence output
	cmd.Stderr = nil
	return cmd.Run()
}
