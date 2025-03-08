package service

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func createServiceFileContents(appDir string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	serviceTemplate := `[Unit]
Description=Audetic Link API
After=network.target

[Service]
User=%s
WorkingDirectory=%s
ExecStart=%s/.venv/bin/gunicorn audetic_link.main:app \
    --bind 0.0.0.0:8481 \
    --worker-class uvicorn.workers.UvicornWorker \
    --log-level debug \
    --log-config %s/src/audetic_link/log/logger.conf
Restart=always
PIDFile=/run/audetic-link.pid

[Install]
WantedBy=multi-user.target`

	return fmt.Sprintf(serviceTemplate, currentUser.Username, appDir, appDir, appDir), nil
}

// Setup installs and starts the systemd service
func Setup(appDir string) error {
	// Copy service file
	contents, err := createServiceFileContents(appDir)
	if err != nil {
		return fmt.Errorf("failed to create service file contents: %w", err)
	}

	dst := "/etc/systemd/system/audetic_link.service"
	if err := os.WriteFile(dst, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// Enable and start service
	if err := exec.Command("systemctl", "enable", "audetic_link.service").Run(); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	if err := exec.Command("systemctl", "start", "audetic_link.service").Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// SetupRemoteAccess configures rpi-connect and user lingering
func SetupRemoteAccess() error {
	// Enable and start rpi-connect
	if err := exec.Command("systemctl", "enable", "rpi-connect").Run(); err != nil {
		return fmt.Errorf("failed to enable rpi-connect: %w", err)
	}

	if err := exec.Command("systemctl", "start", "rpi-connect").Run(); err != nil {
		return fmt.Errorf("failed to start rpi-connect: %w", err)
	}

	// Enable user lingering
	if err := exec.Command("loginctl", "enable-linger", "audetic").Run(); err != nil {
		return fmt.Errorf("failed to enable user lingering: %w", err)
	}

	return nil
}
