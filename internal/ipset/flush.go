package ipset

import (
	"fmt"
	"os/exec"
)

func Flush(setName string) error {
	cmd := exec.Command("ipset", "flush", setName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to flush ipset set '%s': %w", setName, err)
	}
	return nil
}
