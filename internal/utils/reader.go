package utils

import (
	"bufio"
	"io"
	"strings"
)

// ReadEntriesFromReader reads lines from an io.Reader, trims spaces, skips empty lines, and returns a slice of entries.
// Each line is considered a separate entry.
func ReadEntriesFromReader(r io.Reader) ([]string, error) {
	var entries []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}
		entries = append(entries, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}
