package repo

import (
	"fmt"
	"os/exec"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
)

// SetupApp clones the repository and sets up directories
func SetupApp(cfg *types.RPiConfig, repoURL string) error {
	if cfg.DryRun {
		return nil
	}

	// Clone repository
	cfg.Progress.UpdateMessage("Cloning repository...")
	if err := exec.Command("git", "clone", repoURL, cfg.AppDir).Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}
