package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/silvabyte/audeticlinkinstaller/internal/user_utils"
)

const envTemplate = `APP_DIR=%s
AUDETIC_API_URL=https://app.audetic.ai`

// SetupEnv creates and configures the .env file
func SetupEnv(appDir string) error {
	envPath := filepath.Join(appDir, ".env")
	content := fmt.Sprintf(envTemplate, appDir)

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	uid, err := user_utils.UIdToInt(currentUser.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert uid to int: %w", err)
	}

	// Set ownership
	if err := os.Chown(envPath, uid, -1); err != nil {
		return fmt.Errorf("failed to set .env ownership: %w", err)
	}

	return nil
}
