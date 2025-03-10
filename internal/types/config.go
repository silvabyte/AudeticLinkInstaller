package types

import (
	"github.com/silvabyte/AudeticLinkInstaller/internal/pin"
)

// RPiConfig holds configuration for Raspberry Pi installation
type RPiConfig struct {
	ConfigPath   string
	AppDir       string
	RepoUser     string
	RepoOrg      string
	RepoName     string
	RepoToken    string
	ClientID     string
	ClientSecret string
	Progress     *pin.Pin
	DryRun       bool
	Debug        bool
}
