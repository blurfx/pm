package project

import (
	"encoding/json"
	"os"
	"path/filepath"

	"pm/internal/detector"
)

// IsTypeScript checks if the current project uses TypeScript
func IsTypeScript() bool {
	projectRoot, err := detector.FindProjectRoot()
	if err != nil {
		return false
	}

	// Look for tsconfig.json in project root
	tsconfigPath := filepath.Join(projectRoot, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); err == nil {
		return true
	}

	// Look for typescript in package.json dependencies
	packageJSONPath := filepath.Join(projectRoot, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return false
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	if _, ok := pkg.Dependencies["typescript"]; ok {
		return true
	}
	if _, ok := pkg.DevDependencies["typescript"]; ok {
		return true
	}

	return false
}
