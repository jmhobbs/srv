package cmd

import (
	"fmt"
	"path/filepath"
	"syscall"
)

func GetWorkingDirectory(dir string) (string, error) {
	if dir == "" {
		cwd, err := syscall.Getwd()
		if err != nil {
			return "", fmt.Errorf("unable to determine current working directory: %v", err)
		}
		dir = cwd
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("unable to determine absolute path for %q: %v", dir, err)
	}

	return absDir, nil
}
