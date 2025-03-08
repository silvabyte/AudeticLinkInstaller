package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/silvabyte/audeticlinkinstaller/internal/types"
	"github.com/silvabyte/audeticlinkinstaller/internal/user_utils"
)

const envTemplate = `APP_DIR=%s
AUDETIC_API_URL=https://app.audetic.ai`

// SetupEnv creates and configures the .env file
func SetupEnv(cfg *types.RPiConfig) error {
	cfg.Progress.UpdateMessage("Creating .env file...")
	envPath := filepath.Join(cfg.AppDir, ".env")
	content := fmt.Sprintf(envTemplate, cfg.AppDir)

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	realUser, err := user_utils.GetRealUser()
	if err != nil {
		return fmt.Errorf("failed to get real user: %w", err)
	}

	uid, err := user_utils.UIdToInt(realUser.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert uid to int: %w", err)
	}

	// Set ownership
	cfg.Progress.UpdateMessage("Setting file permissions...")
	if err := os.Chown(envPath, uid, -1); err != nil {
		return fmt.Errorf("failed to set .env ownership: %w", err)
	}

	return nil
}
