package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func InitializeDotDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	dotDirPath := filepath.Join(homeDir, ".godo")

	if exists, err := dirExists(dotDirPath); err != nil {
		return fmt.Errorf("error checking dot directory: %w", err)
	} else if !exists {
		if err := os.MkdirAll(dotDirPath, 0755); err != nil {
			return fmt.Errorf("error creating dot directory: %w", err)
		}
	}

	return nil
}
