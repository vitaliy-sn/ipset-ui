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
// Returns the number of successfully added entries, the list of failed entries,
// and a combined error for failed additions.
// If progress is not nil, it is called after each processed entry with the
// number of successfully added entries, the number of processed entries and the total.
func AddEntries(setName string, entries []string, comment string, progress func(added, processed, total int)) (int, []string, error) {
	added := 0
	var failed []string
	var merr *multierror.Error
	total := len(entries)

	for i, entry := range entries {
		if err := addEntry(setName, entry, comment); err != nil {
			failed = append(failed, entry)
			merr = multierror.Append(merr, fmt.Errorf("entry %q: %w", entry, err))
		} else {
			added++
		}
		if progress != nil {
			progress(added, i+1, total)
		}
	}

	err := Save(setName)
	if err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to save set after batch add: %w", err))
	}

	if merr != nil {
		return added, failed, merr.ErrorOrNil()
	}

	return added, failed, nil
}

func deleteEntry(setName, entry string) error {
	cmd := exec.Command("ipset", "del", setName, entry)
	return cmd.Run()
}

func DeleteEntry(setName, entry string) error {
	if err := deleteEntry(setName, entry); err != nil {
		return err
	}

	return Save(setName)
}

// DeleteEntries removes multiple entries from the specified ipset set.
// Returns the number of successfully deleted entries, the list of failed entries,
// and a combined error for failed deletions.
func DeleteEntries(setName string, entries []string) (int, []string, error) {
	deleted := 0
	var failed []string
	var merr *multierror.Error

	for _, entry := range entries {
		if err := deleteEntry(setName, entry); err != nil {
			failed = append(failed, entry)
			merr = multierror.Append(merr, fmt.Errorf("entry %q: %w", entry, err))
			continue
		}
		deleted++
	}

	if err := Save(setName); err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to save set after batch delete: %w", err))
	}

	if merr != nil {
		return deleted, failed, merr.ErrorOrNil()
	}

	return deleted, failed, nil
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
