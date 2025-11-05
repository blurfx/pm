package detector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const maxTraverseDepth = 20

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type packageJSONConfig struct {
	PackageManager string `json:"packageManager"`
}

func detectPackageManager() (PackageManager, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	rootDir, err := FindProjectRoot()
	if err != nil {
		return "", fmt.Errorf("no package manager detected")
	}
	if fileExists(filepath.Join(rootDir, "package-lock.json")) {
		return NPM, nil
	} else if fileExists(filepath.Join(rootDir, "yarn.lock")) {
		return Yarn, nil
	} else if fileExists(filepath.Join(rootDir, "pnpm-lock.yaml")) {
		return Pnpm, nil
	} else if fileExists(filepath.Join(rootDir, "bun.lock")) || fileExists(filepath.Join(currentDir, "bun.lockb")) {
		return Bun, nil
	}

	packageJSONPath, err := FindPackageJSON()
	if err != nil {
		return "", fmt.Errorf("no package manager detected")
	}

	file, err := os.Open(packageJSONPath)
	if err == nil {
		defer file.Close()
		var packageJSON packageJSONConfig
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&packageJSON)
		if err == nil && packageJSON.PackageManager != "" {
			return PackageManager(strings.Split(packageJSON.PackageManager, "@")[0]), nil
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
		return NPM, nil
	}
	if isCommandAvailable("yarn") {
		return Yarn, nil
	}
	if isCommandAvailable("pnpm") {
		return Pnpm, nil
	}
	if isCommandAvailable("bun") {
		return Bun, nil
	}

	return "", fmt.Errorf("no package manager detected (supported: npm, yarn, pnpm, bun)")
}

// Detect detects the package manager used in the current project
func Detect() (PackageManager, error) {
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
		depth++
		if parentDir == currentDir || depth > maxTraverseDepth {
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
		depth++
		if parentDir == currentDir || depth > maxTraverseDepth {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("package.json not found")
}
