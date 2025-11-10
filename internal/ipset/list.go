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
