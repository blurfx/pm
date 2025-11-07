package detector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	packageJSONPath := filepath.Join(rootDir, "package.json")
	packageManagerField := readPackageManagerField(packageJSONPath)

	if fileExists(filepath.Join(rootDir, "package-lock.json")) {
		return NPM, nil
	} else if fileExists(filepath.Join(rootDir, "yarn.lock")) {
		return detectYarnVariant(rootDir, packageManagerField), nil
	} else if fileExists(filepath.Join(rootDir, "pnpm-lock.yaml")) {
		return Pnpm, nil
	} else if fileExists(filepath.Join(rootDir, "bun.lock")) || fileExists(filepath.Join(currentDir, "bun.lockb")) {
		return Bun, nil
	}

	if pm, ok := packageManagerFromField(packageManagerField); ok {
		return pm, nil
	}

	return "", fmt.Errorf("no package manager detected")
}

func readPackageManagerField(packageJSONPath string) string {
	file, err := os.Open(packageJSONPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	var packageJSON packageJSONConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&packageJSON); err != nil {
		return ""
	}

	return packageJSON.PackageManager
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

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func detectYarnVariant(rootDir, packageManagerField string) PackageManager {
	if pm, ok := packageManagerFromField(packageManagerField); ok {
		if pm == Yarn || pm == YarnBerry {
			return pm
		}
	}

	if fileExists(filepath.Join(rootDir, ".yarnrc.yml")) || fileExists(filepath.Join(rootDir, ".yarnrc.yaml")) {
		return YarnBerry
	}
	if fileExists(filepath.Join(rootDir, ".pnp.cjs")) || fileExists(filepath.Join(rootDir, ".pnp.mjs")) {
		return YarnBerry
	}
	if dirExists(filepath.Join(rootDir, ".yarn", "releases")) {
		return YarnBerry
	}
	return Yarn
}

func packageManagerFromField(field string) (PackageManager, bool) {
	name, version := parsePackageManagerParts(field)
	switch name {
	case "":
		return "", false
	case "npm":
		return NPM, true
	case "yarn":
		if isYarnBerryVersion(version) {
			return YarnBerry, true
		}
		return Yarn, true
	case "pnpm":
		return Pnpm, true
	case "bun":
		return Bun, true
	default:
		return PackageManager(name), true
	}
}

func parsePackageManagerParts(value string) (string, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ""
	}

	parts := strings.SplitN(value, "@", 2)
	name := strings.TrimSpace(parts[0])
	version := ""
	if len(parts) == 2 {
		version = strings.TrimSpace(parts[1])
	}

	return name, version
}

func isYarnBerryVersion(version string) bool {
	major, ok := parseSemverMajor(version)
	if !ok {
		return false
	}
	return major >= 2
}

func parseSemverMajor(version string) (int, bool) {
	version = strings.TrimSpace(version)
	if version == "" {
		return 0, false
	}

	version = strings.TrimPrefix(version, "v")
	end := 0
	for end < len(version) && version[end] >= '0' && version[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0, false
	}

	major, err := strconv.Atoi(version[:end])
	if err != nil {
		return 0, false
	}

	return major, true
}
