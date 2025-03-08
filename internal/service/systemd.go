package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
	"github.com/silvabyte/AudeticLinkInstaller/internal/user_utils"
)

func createServiceFileContents(appDir string) (string, error) {
	realUser, err := user_utils.GetRealUser()
	if err != nil {
		return "", fmt.Errorf("failed to get real user: %w", err)
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

	return fmt.Sprintf(serviceTemplate, realUser.Username, appDir, appDir, appDir), nil
}

// Setup installs and starts the systemd service
func Setup(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	// Copy service file
	cfg.Progress.UpdateMessage("Creating service file...")
	contents, err := createServiceFileContents(cfg.AppDir)
	if err != nil {
		return fmt.Errorf("failed to create service file contents: %w", err)
	}

	dst := "/etc/systemd/system/audetic_link.service"
	if err := os.WriteFile(dst, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Set proper ownership of application directory
	realUser, err := user_utils.GetRealUser()
	if err != nil {
		return fmt.Errorf("failed to get real user: %w", err)
	}

	uid, err := user_utils.UIdToInt(realUser.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert uid to int: %w", err)
	}

	cfg.Progress.UpdateMessage("Setting directory permissions...")
	if err := os.Chown(cfg.AppDir, uid, -1); err != nil {
		return fmt.Errorf("failed to set app directory ownership: %w", err)
	}

	// Recursively set ownership of all files
	if err := exec.Command("chown", "-R", fmt.Sprintf("%s:", realUser.Username), cfg.AppDir).Run(); err != nil {
		return fmt.Errorf("failed to set recursive ownership: %w", err)
	}

	// Reload systemd
	cfg.Progress.UpdateMessage("Reloading systemd...")
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// Enable and start service
	cfg.Progress.UpdateMessage("Enabling service...")
	if err := exec.Command("systemctl", "enable", "audetic_link.service").Run(); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	cfg.Progress.UpdateMessage("Starting service...")
	if err := exec.Command("systemctl", "start", "audetic_link.service").Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// SetupRemoteAccess configures rpi-connect and user lingering
func SetupRemoteAccess(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	realUser, err := user_utils.GetRealUser()
	if err != nil {
		return fmt.Errorf("failed to get real user: %w", err)
	}

	// Enable and start rpi-connect
	cfg.Progress.UpdateMessage("Enabling rpi-connect...")
	cmd := exec.Command("systemctl", "enable", "rpi-connect")
	if cfg.Debug {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable rpi-connect: %w\nError output: %s", err, stderr.String())
		}
	} else {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable rpi-connect: %w", err)
		}
	}

	cfg.Progress.UpdateMessage("Starting rpi-connect...")
	cmd = exec.Command("systemctl", "start", "rpi-connect")
	if cfg.Debug {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start rpi-connect: %w\nError output: %s", err, stderr.String())
		}
	} else {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start rpi-connect: %w", err)
		}
	}

	// Enable user lingering
	cfg.Progress.UpdateMessage("Enabling user lingering...")
	cmd = exec.Command("loginctl", "enable-linger", realUser.Username)
	if cfg.Debug {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable user lingering: %w\nError output: %s", err, stderr.String())
		}
	} else {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable user lingering: %w", err)
		}
	}

	return nil
}
