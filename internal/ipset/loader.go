package ipset

import (
	"fmt"
	"ipset-ui/internal/config"
	"ipset-ui/internal/logger"
	"path/filepath"
	"slices"
	"strings"
)

func listSaveFiles() ([]string, error) {
	pattern := filepath.Join(config.Config.AppDir, "*.save")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list save files: %w", err)
	}
	return files, nil
}

func LoadAll() error {
	logger.Info("starting to load all ipset save files...")

	files, err := listSaveFiles()
	if err != nil {
		logger.Error("Error listing save files", "error", err)
		return fmt.Errorf("LoadAll: %w", err)
	}
	logger.Info("found save files", "files", files)

	activeSets, err := List()
	if err != nil {
		logger.Error("Error listing active ipset sets", "error", err)
		return fmt.Errorf("LoadAll: %w", err)
	}
	logger.Info("active ipset sets", "sets", activeSets)

	for _, file := range files {
		base := filepath.Base(file)
		setName := strings.TrimSuffix(base, ".save")
		if slices.Contains(activeSets, setName) {
			logger.Info("flushing ipset set", "set", setName)
			if err := Flush(setName); err != nil {
				logger.Error("error flushing set", "set", setName, "error", err)
				return fmt.Errorf("LoadAll: %w", err)
			}
		}
	}

	for _, file := range files {
		logger.Info("restoring ipset from file", "file", file)
		err := Restore(file)
		if err != nil {
			logger.Error("error restoring from file", "file", file, "error", err)
			return fmt.Errorf("LoadAll: %w", err)
		}
	}

	logger.Info("successfully loaded all ipset save files.")
	return nil
}
