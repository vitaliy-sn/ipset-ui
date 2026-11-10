package ipset

import (
	"os/exec"
	"strings"
)

func List() ([]string, error) {
	cmd := exec.Command("ipset", "list", "-name")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sets := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(sets) == 1 && sets[0] == "" {
		return []string{}, nil
	}
	return sets, nil
}

// SetExists reports whether a set with the given name exists.
func SetExists(setName string) (bool, error) {
	sets, err := List()
	if err != nil {
		return false, err
	}
	for _, s := range sets {
		if s == setName {
			return true, nil
		}
	}
	return false, nil
}
