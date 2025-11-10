package utils

import (
	"fmt"
	"os/exec"
)

// Whois executes the 'whois' command for the given object (IP or domain)
// and returns the output as a string.
func Whois(object string) (string, error) {
	cmd := exec.Command("whois", object)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("whois command failed: %w", err)
	}
	return string(output), nil
}
