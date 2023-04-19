package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const MAX_TRAVERSE_DEPTH = 20

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func detectPackageManager() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}
	depth := 0
	for {
		if fileExists("package-lock.json") {
			return "npm", nil
		} else if fileExists("yarn.lock") {
			return "yarn", nil
		} else if fileExists("pnpm-lock.yaml") {
			return "pnpm", nil
		}

		parentDir := filepath.Dir(currentDir)
		depth += 1
		if parentDir == currentDir || depth > MAX_TRAVERSE_DEPTH {
			break
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("no package manager detected")
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func detectInstalledPackageManagers() (string, error) {
	if isCommandAvailable("npm") {
		return "npm", nil
	}
	if isCommandAvailable("yarn") {
		return "yarn", nil
	}
	if isCommandAvailable("pnpm") {
		return "pnpm", nil
	}

	return "", fmt.Errorf("no package manager detected")
}

func DetectPackageManager() (string, error) {
	packageManager, err := detectPackageManager()
	if err == nil {
		return packageManager, nil
	}
	return detectInstalledPackageManagers()
}
