package ipset

import (
	"os/exec"
)

func Create(setName string) error {
	cmd := exec.Command("ipset", "create", setName, "hash:net", "comment")

	err := cmd.Run()
	if err != nil {
		return err
	}

	return Save(setName)
}
