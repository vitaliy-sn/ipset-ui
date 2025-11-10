package ipset

import (
	"fmt"
	"os/exec"
)

func Restore(path string) error {
	cmd := exec.Command("ipset", "restore", "-exist", "-f", path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore ipset from file '%s': %w", path, err)
	}
	return nil
}
