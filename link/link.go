package link

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/silvabyte/audeticlinkinstaller/internal/installer"
)

type InstallCmd struct {
	Device string `arg:"" help:"Device type to install (rpi02w)" enum:"rpi02w"`
}

type LinkInstaller struct {
	Install     InstallCmd `cmd:"" help:"Install Audetic Link on a device"`
	GithubToken string     `help:"GitHub personal access token for repository access" required:"" env:"GITHUB_TOKEN"`
}

// InstallRPi02W installs Audetic Link for Raspberry Pi Zero 2 W
func InstallRPi02W(cmd LinkInstaller) error {
	info := color.New(color.FgCyan).PrintlnFunc()
	info("Installing Audetic Link for Raspberry Pi Zero 2 W...")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cfg := installer.RPiConfig{
		ConfigPath: "/boot/firmware/config.txt",
		AppDir:     fmt.Sprintf("%s/AudeticLink", home),
		RepoUser:   "matsilva",
		RepoToken:  cmd.GithubToken,
		RepoOrg:    "silvabyte",
		RepoName:   "AudeticLink",
	}

	if err := installer.InstallRPi(cfg); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	return nil
}

// InstallRPiPico installs Audetic Link for Raspberry Pi Pico
func InstallRPiPico() error {
	info := color.New(color.FgCyan).PrintlnFunc()
	info("Installing Audetic Link for Raspberry Pi Pico...")
	info("RPi Pico support coming soon!")
	return nil
}

// Run executes the installer for the specified device
func (cmd LinkInstaller) Run() error {
	switch cmd.Install.Device {
	case "rpi02w":
		return InstallRPi02W(cmd)
	default:
		return fmt.Errorf("unsupported device type: %s", cmd.Install.Device)
	}
}
