package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PackageManager string

const (
	PackageManagerNpm  PackageManager = "npm"
	PackageManagerYarn PackageManager = "yarn"
	PackageManagerPnpm PackageManager = "pnpm"
	PackageManagerBun  PackageManager = "bun"
)

const MAX_TRAVERSE_DEPTH = 20

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type packageJsonConfig struct {
	PackageManager string `json:"packageManager"`
}

func detectPackageManager() (PackageManager, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	rootDir, err := FindProjectRoot()
	if err == nil {
		return "", fmt.Errorf("no package manager detected")
	}
	if fileExists(filepath.Join(rootDir, "package-lock.json")) {
		return PackageManagerNpm, nil
	} else if fileExists(filepath.Join(rootDir, "yarn.lock")) {
		return PackageManagerYarn, nil
	} else if fileExists(filepath.Join(rootDir, "pnpm-lock.yaml")) {
		return PackageManagerPnpm, nil
	} else if fileExists(filepath.Join(rootDir, "bun.lock")) || fileExists(filepath.Join(currentDir, "bun.lockb")) {
		return PackageManagerBun, nil
	}

	packageJsonPath, err := FindPackageJSON()
	if err != nil {
		return "", fmt.Errorf("no package manager detected")
	}

	file, err := os.Open(packageJsonPath)
	if err == nil {
		defer file.Close()
		var packageJson packageJsonConfig
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&packageJson)
		if err == nil && packageJson.PackageManager != "" {
			return PackageManager(strings.Split(packageJson.PackageManager, "@")[0]), nil
		}
	}

	return "", fmt.Errorf("no package manager detected")
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func detectInstalledPackageManagers() (PackageManager, error) {
	if isCommandAvailable("npm") {
		return PackageManagerNpm, nil
	}
	if isCommandAvailable("yarn") {
		return PackageManagerYarn, nil
	}
	if isCommandAvailable("pnpm") {
		return PackageManagerPnpm, nil
	}
	if isCommandAvailable("bun") {
		return PackageManagerBun, nil
	}

	return "", fmt.Errorf("no package manager detected (supported: npm, yarn, pnpm, bun)")
}

func DetectPackageManager() (PackageManager, error) {
	packageManager, err := detectPackageManager()
	if err == nil {
		return packageManager, nil
	}
	return detectInstalledPackageManagers()
}

// FindPackageJSON traverses up the directory tree to find package.json
func FindPackageJSON() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	depth := 0

	for {
		packageJSONPath := filepath.Join(currentDir, "package.json")
		if fileExists(packageJSONPath) {
			return packageJSONPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		depth += 1
		if parentDir == currentDir || depth > MAX_TRAVERSE_DEPTH {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("package.json not found")
}

// FindProjectRoot finds the directory containing package.json by traversing up
func FindProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	depth := 0

	for {
		packageJSONPath := filepath.Join(currentDir, "package.json")
		if fileExists(packageJSONPath) {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		depth += 1
		if parentDir == currentDir || depth > MAX_TRAVERSE_DEPTH {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("package.json not found")
}
