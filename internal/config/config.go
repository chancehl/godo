package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const DotDirectory = ".godo"
const GistFile = "gist_file_id"

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if directory exists: %w", err)
	}
	return info.IsDir(), nil
}

func InitializeDotDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	dotDirPath := filepath.Join(homeDir, DotDirectory)
	exists, err := dirExists(dotDirPath)
	if err != nil {
		return err
	}

	if !exists {
		if err := os.MkdirAll(dotDirPath, 0755); err != nil {
			return fmt.Errorf("error creating dot directory: %w", err)
		}
	}

	return nil
}

func CheckIfGistIdFileExists() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("error getting user home directory: %w", err)
	}

	path := filepath.Join(homeDir, DotDirectory, GistFile)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if gist_id file exists: %w", err)
	}

	return true, nil
}

func WriteGistIdFile(id string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	path := filepath.Join(homeDir, DotDirectory, GistFile)
	if err := os.WriteFile(path, []byte(id), 0644); err != nil {
		return "", fmt.Errorf("error writing gist_id file: %w", err)
	}

	return path, nil
}

func ReadGistIdFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	path := filepath.Join(homeDir, DotDirectory, GistFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading gist_id file: %w", err)
	}

	return string(data), nil
}
