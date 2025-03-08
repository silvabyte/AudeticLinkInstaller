package audio

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/silvabyte/audeticlinkinstaller/internal/types"
)

const alsaConfig = `pcm.i2s_mic {
    type hw
    card sndrpii2scard
    device 0
}

pcm.!default {
    type plug
    slave.pcm "i2s_mic"
}`

// Configure sets up I2S audio on the Raspberry Pi
func Configure(cfg *types.RPiConfig) error {
	// Enable I2S in config.txt
	cfg.Progress.UpdateMessage("Enabling I2S in config.txt...")
	if err := appendIfNotExists(cfg.ConfigPath, "dtparam=i2s=on"); err != nil {
		return fmt.Errorf("failed to enable I2S: %w", err)
	}

	// Add I2S mic overlay
	cfg.Progress.UpdateMessage("Adding I2S mic overlay...")
	if err := appendIfNotExists(cfg.ConfigPath, "dtoverlay=i2s-mmap"); err != nil {
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

	if !strings.Contains(string(content), line) {
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
