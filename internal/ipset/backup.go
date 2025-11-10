package ipset

import (
	"fmt"
	"ipset-ui/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ListBackups(setName string) ([]string, error) {
	fmt.Println("Backup directory:", config.Config.BackupDir)
	pattern := fmt.Sprintf("%s*.save", setName)
	files, err := filepath.Glob(filepath.Join(config.Config.BackupDir, pattern))
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

func DeleteBackup(setName, fileNamePart string) error {
	fileName := fmt.Sprintf("%s-%s.save", setName, fileNamePart)
	filePath := filepath.Join(config.Config.BackupDir, fileName)
	return os.Remove(filePath)
}

func CreateBackup(setName, fileName string) error {
	filePath := filepath.Join(config.Config.BackupDir, fileName)
	cmd := exec.Command("ipset", "save", setName, "-f", filePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create backup for set %q to file %q: %w", setName, filePath, err)
	}
	return nil
}
