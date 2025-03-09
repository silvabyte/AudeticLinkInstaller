package installer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/silvabyte/AudeticLinkInstaller/internal/ascii"
	"github.com/silvabyte/AudeticLinkInstaller/internal/audio"
	"github.com/silvabyte/AudeticLinkInstaller/internal/config"
	"github.com/silvabyte/AudeticLinkInstaller/internal/pin"
	"github.com/silvabyte/AudeticLinkInstaller/internal/python"
	"github.com/silvabyte/AudeticLinkInstaller/internal/repo"
	"github.com/silvabyte/AudeticLinkInstaller/internal/service"
	"github.com/silvabyte/AudeticLinkInstaller/internal/system"
	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
)

const (
	PisugarConfigJSON = `{
	{
    "i2c_bus": 1,
    "double_tap_enable": true,
    "double_tap_shell": "curl -X POST http://0.0.0.0:8481/record/toggle",
    "auto_shutdown_level": 5,
    "auto_shutdown_delay": 30,
    "auto_charging_range": [
        70,
        95
    ],
    "full_charge_duration": 110,
    "auto_power_on": true,
    "soft_poweroff": true,
    "auto_rtc_sync": true
}`
)

// InstallRPi performs the complete installation for Raspberry Pi
func InstallRPi(cfg *types.RPiConfig) error {
	// Check root only for actual installation
	if !cfg.DryRun && os.Geteuid() != 0 {
		return fmt.Errorf("this installer must be run as root for actual installation")
	}

	ascii.Logo()
	fmt.Println("Welcome to the Audetic Link installer for Raspberry Pi!")

	// Initialize progress spinner
	cfg.Progress = pin.New("Starting installation...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorWhite))
	cancel := cfg.Progress.Start(context.Background())
	defer cancel()

	info := color.New(color.FgYellow).PrintlnFunc()
	if cfg.DryRun {
		info("\nDRY RUN: The following changes would be made:")
		info("-------------------------------------------")
		info(fmt.Sprintf("• Application directory: %s", cfg.AppDir))
		info(fmt.Sprintf("• Configuration file: %s", cfg.ConfigPath))
		info("• Service file: /etc/systemd/system/audetic_link.service")
		info("• System packages to be installed:")
		info("  - git")
		info("  - python3-pip")
		info("  - python3-venv")
		info("  - python3-gpiozero")
		info("  - sox")
		info("  - libsox-fmt-all")
		info("  - netcat-openbsd")
		info("  - rpi-connect-lite")
		info("  - i2c-tools")
		info("  - python3-smbus")
		info("  - alsa-utils")
		info("\nConfiguration changes:")
		info(fmt.Sprintf("• I2S will be enabled in %s", cfg.ConfigPath))
		info("• ALSA configuration will be written to /etc/asound.conf")
		info(fmt.Sprintf("• Python virtual environment will be created in %s/.venv", cfg.AppDir))
		info(fmt.Sprintf("• Environment file will be created at %s/.env", cfg.AppDir))
		info("• Systemd service 'audetic_link' will be created and enabled")
		info("• Remote access service 'rpi-connect' will be enabled")
		info("\nNo changes have been made. Run without --dry-run to perform the installation.")
		return nil
	}

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

	//print pisugar install instructions
	// 	wget https://cdn.pisugar.com/release/pisugar-power-manager.sh
	// bash pisugar-power-manager.sh -c release

	fmt.Println("\nPiSugar installation instructions:")
	fmt.Println("wget https://cdn.pisugar.com/release/pisugar-power-manager.sh")
	fmt.Println("bash pisugar-power-manager.sh -c release")
	//read and print json config file
	fmt.Println(PisugarConfigJSON)
	fmt.Println("add the above json to the pisugar config file: /etc/pisugar-server/config.json")

	fmt.Printf("[INFO] Please configure your API credentials in %s/.env\n", cfg.AppDir)
	fmt.Println("[INFO] Reboot your Raspberry Pi to apply all changes: sudo reboot")

	return nil
}
