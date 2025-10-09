package ipset

import (
	"bufio"
	"fmt"
	"ipset-ui/internal/config"
	"ipset-ui/internal/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type EntryWithComment struct {
	Entry   string
	Comment string
}

type IPSet struct{}

func NewIPSet() *IPSet {
	return &IPSet{}
}

func (i *IPSet) SetExists(setName string) (bool, error) {
	cmd := exec.Command("ipset", "list", setName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "does not exist") {
			return false, nil
		}
		return false, fmt.Errorf("error checking if set exists: %v, output: %s", err, string(output))
	}
	return true, nil
}

func (i *IPSet) CreateSet(setName string) error {
	cmd := exec.Command("ipset", "create", setName, "hash:net", "comment")
	return cmd.Run()
}

func (i *IPSet) DeleteSet(setName string) error {
	cmd := exec.Command("ipset", "destroy", setName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("set not deleted: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (i *IPSet) ListSets() ([]string, error) {
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

func (i *IPSet) AddEntry(setName, entry, comment string) error {
	cmd := exec.Command("ipset", "add", setName, entry, "comment", comment)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("entry not added: %s, error: %v, output: %s", entry, err, strings.TrimSpace(string(output)))
	}
	return nil
}

func (i *IPSet) DeleteEntry(setName, entry string) error {
	cmd := exec.Command("ipset", "del", setName, entry)
	return cmd.Run()
}

func (i *IPSet) ListEntries(setName string) ([]EntryWithComment, error) {
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

// AddEntries adds multiple entries to the specified ipset set.
// Returns the number of successfully added entries and error if any critical error occurs.
func (i *IPSet) AddEntries(setName string, entries []string, comment string) (int, error) {
	added := 0
	for _, entry := range entries {
		err := i.AddEntry(setName, entry, comment)
		if err == nil {
			added++
			continue
		}

		// For other errors, return immediately
		if err != nil {
			return added, err
		}
	}
	return added, nil
}

func (i *IPSet) SaveSet(setName, fileName string) error {
	filePath := filepath.Join(config.AppConfig.IpsetBackupDir, fileName)
	cmd := exec.Command("ipset", "save", setName, "-f", filePath)
	return cmd.Run()
}

func (i *IPSet) RestoreSet(setName, fileNamePart string) error {
	err := i.DeleteSet(setName)
	if err != nil {
		return fmt.Errorf("error deleting set: %v", err)
	}

	fileName := fmt.Sprintf("%s-%s.save", setName, fileNamePart)
	filePath := filepath.Join(config.AppConfig.IpsetBackupDir, fileName)
	log.Println("filePath: ", filePath)
	cmd := exec.Command("ipset", "restore", "-f", filePath)
	log.Println("cmd: ", cmd.Stdout, cmd.Stderr)
	return cmd.Run()
}

func (i *IPSet) ListBackupFiles(setName string) ([]string, error) {
	pattern := fmt.Sprintf("%s*.save", setName)
	files, err := filepath.Glob(filepath.Join(config.AppConfig.IpsetBackupDir, pattern))
	if err != nil {
		return nil, err
	}

	var backupFiles []string
	for _, file := range files {
		baseName := filepath.Base(file)
		trimmedName := strings.TrimPrefix(baseName, setName+"-")
		trimmedName = strings.TrimSuffix(trimmedName, ".save")
		backupFiles = append(backupFiles, trimmedName)
	}

	return backupFiles, nil
}

func (i *IPSet) DeleteBackupFile(setName, fileNamePart string) error {
	fileName := fmt.Sprintf("%s-%s.save", setName, fileNamePart)
	filePath := filepath.Join(config.AppConfig.IpsetBackupDir, fileName)
	return os.Remove(filePath)
}
