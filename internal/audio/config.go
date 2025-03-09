package audio

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/silvabyte/AudeticLinkInstaller/internal/types"
)

const alsaConfig = `pcm.!default {
    type hw
    card 0
}

ctl.!default {
    type hw
    card 0
}`

// Configure sets up I2S audio on the Raspberry Pi
func Configure(cfg *types.RPiConfig) error {
	if cfg.DryRun {
		return nil
	}

	// Enable I2S in config.txt
	// 	dtparam=i2s=on
	// dtoverlay=i2s-mmap
	// dtoverlay=rpi-i2s-mmap
	// dtoverlay=googlevoicehat-soundcard
	cfg.Progress.UpdateMessage("Enabling dtparam=i2s=on")
	if err := appendIfNotExists(cfg.ConfigPath, "dtparam=i2s=on"); err != nil {
		return fmt.Errorf("failed to enable I2S: %w", err)
	}

	// Add I2S mic overlay
	cfg.Progress.UpdateMessage("Adding dtoverlay=i2s-mmap")
	if err := appendIfNotExists(cfg.ConfigPath, "dtoverlay=i2s-mmap"); err != nil {
		return fmt.Errorf("failed to add I2S overlay: %w", err)
	}

	cfg.Progress.UpdateMessage("Adding dtoverlay=rpi-i2s-mmap")
	if err := appendIfNotExists(cfg.ConfigPath, "dtoverlay=rpi-i2s-mmap"); err != nil {
		return fmt.Errorf("failed to add I2S overlay: %w", err)
	}

	cfg.Progress.UpdateMessage("Adding dtoverlay=googlevoicehat-soundcard")
	if err := appendIfNotExists(cfg.ConfigPath, "dtoverlay=googlevoicehat-soundcard"); err != nil {
		return fmt.Errorf("failed to add I2S overlay: %w", err)
	}

	// Write ALSA config
	cfg.Progress.UpdateMessage("Writing ALSA configuration...")
	if err := os.WriteFile("/etc/asound.conf", []byte(alsaConfig), 0644); err != nil {
		return fmt.Errorf("failed to write ALSA config: %w", err)
	}

	// Reload ALSA
	cfg.Progress.UpdateMessage("Reloading ALSA...")
	if err := exec.Command("alsactl", "restore").Run(); err != nil {
		return fmt.Errorf("failed to reload ALSA: %w", err)
	}

	return nil
}

func appendIfNotExists(path, line string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Split content into lines and check each uncommented line
	lines := strings.Split(string(content), "\n")
	exists := false
	for _, existingLine := range lines {
		trimmed := strings.TrimSpace(existingLine)
		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed == strings.TrimSpace(line) {
			exists = true
			break
		}
	}

	if !exists {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString("\n" + line); err != nil {
			return err
		}
	}

	return nil
}
