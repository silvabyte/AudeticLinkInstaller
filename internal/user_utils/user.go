package user_utils

import (
	"fmt"
	"os"
	"os/user"
)

func UIdToInt(uid string) (int, error) {
	var uidInt int
	if _, err := fmt.Sscanf(uid, "%d", &uidInt); err != nil {
		return -1, fmt.Errorf("failed to convert uid to int: %w", err)
	}
	return uidInt, nil
}

// GetRealUser gets the actual user who invoked sudo
func GetRealUser() (*user.User, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser == "" {
		// If not running under sudo, return current user
		return user.Current()
	}
	return user.Lookup(sudoUser)
}

// GetRealUserHome gets the home directory of the actual user who invoked sudo
func GetRealUserHome() (string, error) {
	u, err := GetRealUser()
	if err != nil {
		return "", fmt.Errorf("failed to get real user: %w", err)
	}
	return u.HomeDir, nil
}
