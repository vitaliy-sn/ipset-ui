package ipset

import (
	"fmt"
	"os/exec"

	"ipset-ui/internal/logger"
)

func Create(setName string) error {
	logger.Info("Executing", "command", fmt.Sprintf("ipset create %s hash:net comment", setName))
	cmd := exec.Command("ipset", "create", setName, "hash:net", "comment")

	err := cmd.Run()
	if err != nil {
		return err
	}

	return Save(setName)
}
