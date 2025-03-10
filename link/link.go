package link

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/silvabyte/AudeticLinkInstaller/internal/installer"
	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
	"github.com/silvabyte/AudeticLinkInstaller/internal/user_utils"
)

type InstallCmd struct {
	Device string `arg:"" help:"Device type to install (rpi02w)" enum:"rpi02w"`
	DryRun bool   `help:"Simulate installation without making any changes" default:"false"`
	Debug  bool   `help:"Enable debug mode" default:"false"`
}

type LinkInstaller struct {
	Install      InstallCmd `cmd:"" help:"Install Audetic Link on a device"`
	GithubToken  string     `help:"GitHub personal access token for repository access" required:"" env:"GITHUB_TOKEN"`
	ClientID     string     `help:"Audetic API Client ID for OAuth2 authentication" optional:"" env:"CLIENT_ID"`
	ClientSecret string     `help:"Audetic API Client Secret for OAuth2 authentication" optional:"" env:"CLIENT_SECRET"`
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
		ConfigPath:   "/boot/firmware/config.txt",
		AppDir:       fmt.Sprintf("%s/AudeticLink", home),
		RepoUser:     "matsilva",
		RepoToken:    cmd.GithubToken,
		RepoOrg:      "silvabyte",
		RepoName:     "AudeticLink",
		ClientID:     cmd.ClientID,
		ClientSecret: cmd.ClientSecret,
		DryRun:       cmd.Install.DryRun,
		Debug:        cmd.Install.Debug,
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
