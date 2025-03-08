package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const envTemplate = `APP_DIR=%s
AUDETIC_API_URL=https://app.audetic.ai
LINK_CLIENT_ID=
LINK_CLIENT_SECRET=`

// SetupEnv creates and configures the .env file
func SetupEnv(appDir string) error {
	envPath := filepath.Join(appDir, ".env")
	content := fmt.Sprintf(envTemplate, appDir)

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	// Set ownership
	uid := getUserID("audetic")
	if uid != -1 {
		if err := os.Chown(envPath, uid, -1); err != nil {
			return fmt.Errorf("failed to set .env ownership: %w", err)
		}
	}

	return nil
}

func getUserID(username string) int {
	u, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	uid := -1
	fmt.Sscanf(u.Uid, "%d", &uid)
	return uid
}
