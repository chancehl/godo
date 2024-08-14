package config

import (
	"errors"
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

func CheckIfGistIdFileExists() bool {
	homeDir, _ := os.UserHomeDir()

	path := filepath.Join(homeDir, ".godo", "gist_file_id")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func WriteGistIdFile(id string) (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("error getting user home directory: %s", err)
		return "", err
	}

	path := filepath.Join(homeDir, ".godo", "gist_file_id")
	os.WriteFile(path, []byte(id), 0644)

	return path, nil
}

func ReadGistIdFile() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("error getting user home directory: %s", err)
		return "", err
	}

	path := filepath.Join(homeDir, ".godo", "gist_file_id")

	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("error reading gist_file_id file: %s", err)
		return "", err
	}

	return string(data), nil
}
