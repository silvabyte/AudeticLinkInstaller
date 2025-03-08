package system

import (
	"fmt"
	"os/exec"

	"github.com/silvabyte/audeticlinkinstaller/internal/types"
)

// InstallDependencies installs all required system packages
func InstallDependencies(cfg *types.RPiConfig) error {
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
	cfg.Progress.UpdateMessage("Updating package list...")
	if err := execCommand("apt", "update"); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Upgrade system
	cfg.Progress.UpdateMessage("Upgrading system packages...")
	if err := execCommand("apt", "full-upgrade", "-y"); err != nil {
		return fmt.Errorf("failed to upgrade system: %w", err)
	}

	// Install packages
	for _, dep := range deps {
		cfg.Progress.UpdateMessage(fmt.Sprintf("Installing %s...", dep))
		if err := execCommand("apt", "install", "-y", dep); err != nil {
			return fmt.Errorf("failed to install %s: %w", dep, err)
		}
	}

	return nil
}

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil // Silence output
	cmd.Stderr = nil
	return cmd.Run()
}
