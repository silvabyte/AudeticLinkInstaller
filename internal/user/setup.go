package user

import (
	"fmt"
	"os/exec"
)

// SetupApp clones the repository and sets up directories
func SetupApp(appDir, repoURL string) error {
	// Clone repository
	if err := exec.Command("git", "clone", repoURL, appDir).Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}
