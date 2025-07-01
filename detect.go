package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const MAX_TRAVERSE_DEPTH = 20

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type packageJsonConfig struct {
	PackageManager string `json:"packageManager"`
}

func detectPackageManager() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	depth := 0

	for {
		if fileExists(filepath.Join(currentDir, "package-lock.json")) {
			return "npm", nil
		} else if fileExists(filepath.Join(currentDir, "yarn.lock")) {
			return "yarn", nil
		} else if fileExists(filepath.Join(currentDir, "pnpm-lock.yaml")) {
			return "pnpm", nil
		}

		packageJsonPath := filepath.Join(currentDir, "package.json")
		if fileExists(packageJsonPath) {
			file, err := os.Open(packageJsonPath)
			if err == nil {
				defer file.Close()
				var packageJson packageJsonConfig
				decoder := json.NewDecoder(file)
				err = decoder.Decode(&packageJson)
				if err == nil && packageJson.PackageManager != "" {
					return strings.Split(packageJson.PackageManager, "@")[0], nil
				}
			}
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

	return "", fmt.Errorf("no package manager detected (supported: npm, yarn, pnpm)")
}

func DetectPackageManager() (string, error) {
	packageManager, err := detectPackageManager()
	if err == nil {
		return packageManager, nil
	}
	return detectInstalledPackageManagers()
}
