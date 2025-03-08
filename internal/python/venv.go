package python

import (
	"fmt"
	"os/exec"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
)

// SetupVenv creates and configures the Python virtual environment
func SetupVenv(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	// Create virtual environment
	cfg.Progress.UpdateMessage("Creating Python virtual environment...")
	cmd := exec.Command("python3", "-m", "venv", ".venv")
	cmd.Dir = cfg.AppDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Install packages
	packages := []string{
		"--upgrade", "pip",
		"uv",
		"uvicorn",
		"gunicorn",
		"fastapi",
		"python-dotenv",
		"structlog",
	}

	// Install packages
	for _, pkg := range packages {
		cfg.Progress.UpdateMessage(fmt.Sprintf("Installing Python package: %s", pkg))
		args := fmt.Sprintf("source %s/.venv/bin/activate && pip install %s", cfg.AppDir, pkg)
		cmd = exec.Command("bash", "-c", args)
		cmd.Dir = cfg.AppDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Python package %s: %w", pkg, err)
		}
	}

	// Sync packages
	cfg.Progress.UpdateMessage("Syncing Python packages...")
	args := fmt.Sprintf("source %s/.venv/bin/activate && uv sync", cfg.AppDir)
	cmd = exec.Command("bash", "-c", args)
	cmd.Dir = cfg.AppDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to sync Python packages: %w", err)
	}

	return nil
}
