package installer

import (
	"fmt"
	"os"

	"github.com/silvabyte/audeticlinkinstaller/internal/audio"
	"github.com/silvabyte/audeticlinkinstaller/internal/config"
	"github.com/silvabyte/audeticlinkinstaller/internal/python"
	"github.com/silvabyte/audeticlinkinstaller/internal/service"
	"github.com/silvabyte/audeticlinkinstaller/internal/system"
	"github.com/silvabyte/audeticlinkinstaller/internal/user"
)

// RPiConfig holds configuration for Raspberry Pi installation
type RPiConfig struct {
	ConfigPath string
	AppDir     string
	RepoUser   string
	RepoToken  string
	RepoOrg    string
	RepoName   string
}

// InstallRPi performs the complete installation for Raspberry Pi
func InstallRPi(cfg RPiConfig) error {
	// Check root
	if os.Geteuid() != 0 {
		return fmt.Errorf("this installer must be run as root")
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Installing system dependencies", system.InstallDependencies},
		{"Configuring audio", func() error { return audio.Configure(cfg.ConfigPath) }},
		// {"Setting up user", user.Setup},
		{"Setting up application", func() error {
			return user.SetupApp(cfg.AppDir, fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", cfg.RepoUser, cfg.RepoToken, cfg.RepoOrg, cfg.RepoName))
		}},
		{"Setting up Python environment", func() error { return python.SetupVenv(cfg.AppDir) }},
		{"Setting up environment variables", func() error { return config.SetupEnv(cfg.AppDir) }},
		{"Setting up service", func() error { return service.Setup(cfg.AppDir) }},
		{"Setting up remote access", func() error { return service.SetupRemoteAccess() }},
	}

	for _, step := range steps {
		fmt.Printf("[INFO] %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s failed: %w", step.name, err)
		}
	}

	fmt.Println("[INFO] Installation complete!")
	fmt.Printf("[INFO] Please configure your API credentials in %s/.env\n", cfg.AppDir)
	fmt.Println("[INFO] Reboot your Raspberry Pi to apply all changes")

	return nil
}
