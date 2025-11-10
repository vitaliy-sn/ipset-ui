package ipset

import (
	"fmt"
	"ipset-ui/internal/config"
	"os/exec"
	"path/filepath"
)

func Save(setName string) error {
	filePath := filepath.Join(config.Config.AppDir, setName+".save")
	cmd := exec.Command("ipset", "save", setName, "-f", filePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ipset save failed: %v", err)
	}
	return nil
}
