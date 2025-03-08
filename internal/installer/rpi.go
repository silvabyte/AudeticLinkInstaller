package installer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/silvabyte/audeticlinkinstaller/internal/audio"
	"github.com/silvabyte/audeticlinkinstaller/internal/config"
	"github.com/silvabyte/audeticlinkinstaller/internal/pin"
	"github.com/silvabyte/audeticlinkinstaller/internal/python"
	"github.com/silvabyte/audeticlinkinstaller/internal/repo"
	"github.com/silvabyte/audeticlinkinstaller/internal/service"
	"github.com/silvabyte/audeticlinkinstaller/internal/system"
	"github.com/silvabyte/audeticlinkinstaller/internal/types"
)

// InstallRPi performs the complete installation for Raspberry Pi
func InstallRPi(cfg *types.RPiConfig) error {
	// Check root
	if os.Geteuid() != 0 {
		return fmt.Errorf("this installer must be run as root")
	}

	// Initialize progress spinner
	cfg.Progress = pin.New("Starting installation...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorWhite))
	cancel := cfg.Progress.Start(context.Background())
	defer cancel()

	// Step 1: Install system dependencies
	cfg.Progress.UpdateMessage("Installing system dependencies...")
	if err := system.InstallDependencies(cfg); err != nil {
		cfg.Progress.Fail("System dependencies installation failed")
		return fmt.Errorf("installing system dependencies: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 2: Configure audio
	cfg.Progress.UpdateMessage("Configuring audio...")
	if err := audio.Configure(cfg); err != nil {
		cfg.Progress.Fail("Audio configuration failed")
		return fmt.Errorf("configuring audio: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 3: Setup application
	cfg.Progress.UpdateMessage("Setting up application...")
	repoURL := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git",
		cfg.RepoUser, cfg.RepoToken, cfg.RepoOrg, cfg.RepoName)
	if err := repo.SetupApp(cfg, repoURL); err != nil {
		cfg.Progress.Fail("Application setup failed")
		return fmt.Errorf("setting up application: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 4: Setup Python environment
	cfg.Progress.UpdateMessage("Setting up Python environment...")
	if err := python.SetupVenv(cfg); err != nil {
		cfg.Progress.Fail("Python environment setup failed")
		return fmt.Errorf("setting up Python environment: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 5: Setup environment variables
	cfg.Progress.UpdateMessage("Setting up environment variables...")
	if err := config.SetupEnv(cfg); err != nil {
		cfg.Progress.Fail("Environment setup failed")
		return fmt.Errorf("setting up environment variables: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 6: Setup service
	cfg.Progress.UpdateMessage("Setting up service...")
	if err := service.Setup(cfg); err != nil {
		cfg.Progress.Fail("Service setup failed")
		return fmt.Errorf("setting up service: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Step 7: Setup remote access
	cfg.Progress.UpdateMessage("Setting up remote access...")
	if err := service.SetupRemoteAccess(cfg); err != nil {
		cfg.Progress.Fail("Remote access setup failed")
		return fmt.Errorf("setting up remote access: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	cfg.Progress.Stop("Installation complete!")
	fmt.Printf("[INFO] Please configure your API credentials in %s/.env\n", cfg.AppDir)
	fmt.Println("[INFO] Reboot your Raspberry Pi to apply all changes")

	return nil
}
