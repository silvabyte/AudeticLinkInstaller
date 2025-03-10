package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
	"github.com/silvabyte/AudeticLinkInstaller/internal/user_utils"
)

const envTemplate = `APP_DIR=%s
AUDETIC_API_URL=https://app.audetic.ai%s%s`

// SetupEnv creates and configures the .env file
func SetupEnv(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	cfg.Progress.UpdateMessage("Creating .env file...")
	envPath := filepath.Join(cfg.AppDir, ".env")

	// Prepare optional environment variables
	var clientIDEnv, clientSecretEnv string
	if cfg.ClientID != "" {
		clientIDEnv = fmt.Sprintf("\nLINK_CLIENT_ID=%s", cfg.ClientID)
	}
	if cfg.ClientSecret != "" {
		clientSecretEnv = fmt.Sprintf("\nLINK_CLIENT_SECRET=%s", cfg.ClientSecret)
	}

	content := fmt.Sprintf(envTemplate, cfg.AppDir, clientIDEnv, clientSecretEnv)

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
