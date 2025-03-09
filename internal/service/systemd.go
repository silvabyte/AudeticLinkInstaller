package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
	"github.com/silvabyte/AudeticLinkInstaller/internal/user_utils"
)

// execWithLogging executes a command and logs output to a file
func execWithLogging(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Write outputs to log file
	logPath := "audetic_link_debug.log"
	logContent := fmt.Sprintf("\n=== Command: %v ===\nTimestamp: %s\n\nStdout:\n%s\n\nStderr:\n%s\n\n",
		cmd.Args, time.Now().Format(time.RFC3339), stdout.String(), stderr.String())

	// Open file in append mode or create if doesn't exist
	f, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return fmt.Errorf("failed to open log file: %w", openErr)
	}
	defer f.Close()

	if _, writeErr := f.WriteString(logContent); writeErr != nil {
		return fmt.Errorf("failed to write to log file: %w", writeErr)
	}

	if err != nil {
		return fmt.Errorf("%w\nCheck %s for details", err, logPath)
	}

	// Don't remove the log file since we want to keep the history
	return nil
}

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
	cmd := exec.Command("chown", "-R", fmt.Sprintf("%s:", realUser.Username), cfg.AppDir)
	if err := execWithLogging(cmd); err != nil {
		return fmt.Errorf("failed to set recursive ownership: %w", err)
	}

	// Reload systemd
	cfg.Progress.UpdateMessage("Reloading systemd...")
	cmd = exec.Command("systemctl", "daemon-reload")
	if err := execWithLogging(cmd); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// Enable and start service
	cfg.Progress.UpdateMessage("Enabling service...")
	cmd = exec.Command("systemctl", "enable", "audetic_link.service")
	if err := execWithLogging(cmd); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	cfg.Progress.UpdateMessage("Starting service...")
	cmd = exec.Command("systemctl", "start", "audetic_link.service")
	if err := execWithLogging(cmd); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// SetupRemoteAccess configures rpi-connect and user lingering
func SetupRemoteAccess(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	// Enable and start rpi-connect
	cfg.Progress.UpdateMessage("Enabling rpi-connect...")
	cmd := exec.Command("rpi-connect", "on")
	if err := execWithLogging(cmd); err != nil {
		return fmt.Errorf("failed to enable rpi-connect: %w", err)
	}

	fmt.Println("[INFO] rpi-connect is on, run the following command fully step:")

	fmt.Println("rpi-connect signin")
	fmt.Println("loginctl enable-linger")

	return nil
}
