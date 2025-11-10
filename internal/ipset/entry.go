package ipset

import (
	"bufio"
	"fmt"
	"ipset-ui/internal/utils"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type EntryWithComment struct {
	Entry   string
	Comment string
}

func addEntry(setName, entry, comment string) error {
	cmd := exec.Command("ipset", "add", setName, entry, "comment", comment)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add entry: %s, error: %v, output: %s", entry, err, strings.TrimSpace(string(output)))
	}
	return nil
}

// AddEntries adds multiple entries to the specified ipset set.
// Returns the number of successfully added entries and a combined error for failed additions.
func AddEntries(setName string, entries []string, comment string) (int, error) {
	added := 0
	var merr *multierror.Error

	for _, entry := range entries {
		if err := addEntry(setName, entry, comment); err != nil {
			merr = multierror.Append(merr, fmt.Errorf("entry %q: %w", entry, err))
			continue
		}
		added++
	}

	err := Save(setName)
	if err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to save set after batch add: %w", err))
	}

	if merr != nil {
		return added, merr.ErrorOrNil()
	}

	return added, nil
}

func DeleteEntry(setName, entry string) error {
	cmd := exec.Command("ipset", "del", setName, entry)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return Save(setName)
}

func ListEntries(setName string) ([]EntryWithComment, error) {
	cmd := exec.Command("ipset", "list", setName)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(output)
	var entries []EntryWithComment
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Split by " comment "
		parts := strings.SplitN(line, " comment ", 2)
		entry := strings.TrimSpace(parts[0])
		if !utils.IsValidCIDR(entry) && !utils.IsValidIP(entry) {
			continue
		}
		comment := ""
		if len(parts) == 2 {
			// Extract text in quotes
			commentStart := strings.Index(parts[1], "\"")
			commentEnd := strings.LastIndex(parts[1], "\"")
			if commentStart != -1 && commentEnd > commentStart {
				comment = parts[1][commentStart+1 : commentEnd]
			}
		}
		entries = append(entries, EntryWithComment{Entry: entry, Comment: comment})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return entries, nil
}
