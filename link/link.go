package link

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/silvabyte/audeticlinkinstaller/internal/installer"
	"github.com/silvabyte/audeticlinkinstaller/internal/types"
	"github.com/silvabyte/audeticlinkinstaller/internal/user_utils"
)

type InstallCmd struct {
	Device string `arg:"" help:"Device type to install (rpi02w)" enum:"rpi02w"`
	DryRun bool   `help:"Simulate installation without making any changes" default:"false"`
}

type LinkInstaller struct {
	Install     InstallCmd `cmd:"" help:"Install Audetic Link on a device"`
	GithubToken string     `help:"GitHub personal access token for repository access" required:"" env:"GITHUB_TOKEN"`
}

// InstallRPi02W installs Audetic Link for Raspberry Pi Zero 2 W
func InstallRPi02W(cmd *LinkInstaller) error {
	info := color.New(color.FgCyan).PrintlnFunc()
	if cmd.Install.DryRun {
		info("Simulating installation of Audetic Link for Raspberry Pi Zero 2 W...")
	} else {
		info("Installing Audetic Link for Raspberry Pi Zero 2 W...")
	}

	// Get the real user's home directory
	home, err := user_utils.GetRealUserHome()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	cfg := &types.RPiConfig{
		ConfigPath: "/boot/firmware/config.txt",
		AppDir:     fmt.Sprintf("%s/AudeticLink", home),
		RepoUser:   "matsilva",
		RepoToken:  cmd.GithubToken,
		RepoOrg:    "silvabyte",
		RepoName:   "AudeticLink",
		DryRun:     cmd.Install.DryRun,
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
func (cmd *LinkInstaller) Run() error {
	switch cmd.Install.Device {
	case "rpi02w":
		return InstallRPi02W(cmd)
	default:
		return fmt.Errorf("unsupported device type: %s", cmd.Install.Device)
	}
}
