package python

import (
	"fmt"
	"os/exec"
)

// SetupVenv creates and configures the Python virtual environment
func SetupVenv(appDir string) error {
	// Create virtual environment
	cmd := exec.Command("python3", "-m", "venv", ".venv")
	cmd.Dir = appDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Install packages
	packages := []string{
		"--upgrade", "pip",
		"uv",
		"uvicorn",
		"gunicorn",
		"fastapi",
		"python-dotenv",
		"structlog",
	}

	args := fmt.Sprintf("source %s/.venv/bin/activate && pip install %s", appDir, joinArgs(packages))
	cmd = exec.Command("bash", "-c", args)
	cmd.Dir = appDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Python packages: %w", err)
	}

	args = fmt.Sprintf("source %s/.venv/bin/activate && uv sync", appDir)
	cmd = exec.Command("bash", "-c", args)
	cmd.Dir = appDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to sync Python packages: %w", err)
	}

	return nil
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}
