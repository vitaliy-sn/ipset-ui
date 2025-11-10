package ipset

import (
	"fmt"
	"os/exec"
	"strings"
)

func Destroy(setName string) error {
	cmd := exec.Command("ipset", "destroy", setName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to destroy the set: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
